package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
	amqp "github.com/rabbitmq/amqp091-go"
	resolutionServiceProto "github.com/szewczukk/resolution-service/proto"
	scoreServiceProto "github.com/szewczukk/score-service/proto"
	userServiceProto "github.com/szewczukk/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateServicePayload struct {
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationTokenPayload struct {
	Token string `json:"token"`
}

type CompleteResolutionPayload struct {
	UserId       int32 `json:"userId"`
	ResolutionId int32 `json:"resolutionId"`
}

type GetUserScoresPayload struct {
	UserScores []UserScore
}

type UserScore struct {
	UserId   int32  `json:"userId"`
	Username string `json:"username"`
	Score    int32  `json:"score"`
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func main() {
	scoreServiceConnection, err := grpc.Dial(
		":3002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(err)
	}

	resolutionServiceConnection, err := grpc.Dial(
		":3001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(err)
	}

	userServiceConnection, err := grpc.Dial(
		":3000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(err)
	}

	resolutionServiceClient := resolutionServiceProto.NewResolutionServiceClient(
		resolutionServiceConnection,
	)
	userServiceClient := userServiceProto.NewUserServiceClient(
		userServiceConnection,
	)
	scoreServiceCLient := scoreServiceProto.NewScoreServiceClient(
		scoreServiceConnection,
	)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"resolutionCompleted",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/users/", func(c *fiber.Ctx) error {
		userScores, err := scoreServiceCLient.GetAllUserScores(
			context.Background(),
			&scoreServiceProto.GetAllUserScoresRequest{},
		)

		if err != nil {
			return err
		}

		payload := new(GetUserScoresPayload)

		for _, userScore := range userScores.Scores {
			userId := int32(userScore.UserId)
			user, err := userServiceClient.GetUserById(
				context.Background(),
				&userServiceProto.UserServiceUserId{UserId: int32(userScore.UserId)},
			)
			if err != nil {
				return err
			}

			payload.UserScores = append(payload.UserScores, UserScore{
				UserId:   userId,
				Username: user.Username,
				Score:    userScore.Score,
			})
		}

		return c.JSON(payload)
	})

	app.Get("/resolutions/", func(c *fiber.Ctx) error {
		resolutions, err := resolutionServiceClient.GetAllResolutions(
			context.Background(),
			&resolutionServiceProto.GetAllResolutionsRequest{},
		)

		if err != nil {
			return err
		}

		return c.JSON(resolutions.Resolutions)
	})

	app.Get("/current-user/resolutions/", func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")

		userId, err := getUserIdFromAuthorizationHeader(authorizationHeader)
		if err != nil {
			return err
		}

		resolutions, err := resolutionServiceClient.GetResolutionsByUserId(
			context.Background(),
			&resolutionServiceProto.UserId{UserId: userId},
		)

		if err != nil {
			return err
		}

		return c.JSON(resolutions.Resolutions)
	})

	app.Post("/current-user/resolutions/:id/complete", func(c *fiber.Ctx) error {
		payload := new(CompleteResolutionPayload)

		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return err
		}

		userId, err := getUserIdFromAuthorizationHeader(c.Get("Authorization"))
		if err != nil {
			return err
		}

		payload.ResolutionId = int32(id)
		payload.UserId = userId

		serialized, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = ch.PublishWithContext(
			ctx,
			"",
			queue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        serialized,
			},
		)

		if err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	app.Get("/current-user/", func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")

		userId, err := getUserIdFromAuthorizationHeader(authorizationHeader)
		if err != nil {
			return err
		}

		user, err := userServiceClient.GetUserById(
			context.Background(),
			&userServiceProto.UserServiceUserId{UserId: userId},
		)

		if err != nil {
			return err
		}

		return c.JSON(user)
	})

	app.Post("/resolutions/", func(c *fiber.Ctx) error {
		createResolutionPayload := new(CreateServicePayload)
		if err := c.BodyParser(createResolutionPayload); err != nil {
			return err
		}

		resolution, err := resolutionServiceClient.CreateResolution(
			context.Background(),
			&resolutionServiceProto.CreateResolutionRequest{
				Name:   createResolutionPayload.Name,
				UserId: int32(createResolutionPayload.UserId),
			},
		)

		if err != nil {
			return err
		}

		return c.JSON(resolution)
	})

	app.Post("/login/", func(c *fiber.Ctx) error {
		loginPayload := new(LoginPayload)
		if err := c.BodyParser(loginPayload); err != nil {
			return err
		}

		response, err := userServiceClient.AuthenticateUser(
			context.Background(),
			&userServiceProto.UserCredentials{
				Username: loginPayload.Username,
				Password: loginPayload.Password,
			},
		)
		if err != nil {
			return err
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": response.UserId,
		})

		tokenString, err := token.SignedString(sampleSecretKey)
		if err != nil {
			return err
		}

		return c.JSON(AuthenticationTokenPayload{Token: tokenString})
	})

	app.Listen(":3003")
}

func getUserIdFromAuthorizationHeader(header string) (int32, error) {
	splitToken := strings.Split(header, "Bearer ")
	jwtTokenString := splitToken[1]

	token, err := jwt.Parse(jwtTokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token singing method")
		}

		return sampleSecretKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claim")
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	userId := int32(claims["userId"].(float64))

	return userId, nil
}
