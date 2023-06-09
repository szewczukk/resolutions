package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
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

type GetUserScoresPayload struct {
	UserScores []UserScore
}

type UserScore struct {
	UserId   int32  `json:"userId"`
	Username string `json:"username"`
	Score    int32  `json:"score"`
}

type UserCredentails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CurrentUserCreateResolutionRequest struct {
	Name string `json:"name"`
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func main() {
	scoreServiceConnection, err := grpc.Dial(
		":3003",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(err)
	}

	resolutionServiceConnection, err := grpc.Dial(
		":3002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println(err)
	}

	userServiceConnection, err := grpc.Dial(
		":3001",
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

		return c.JSON(payload.UserScores)
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

	app.Post("/current-user/resolutions/", func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")

		userId, err := getUserIdFromAuthorizationHeader(authorizationHeader)
		if err != nil {
			return err
		}

		body := new(CurrentUserCreateResolutionRequest)
		if err = c.BodyParser(&body); err != nil {
			return err
		}

		resolution, err := resolutionServiceClient.CreateResolution(
			context.Background(),
			&resolutionServiceProto.CreateResolutionRequest{
				UserId: userId,
				Name:   body.Name,
			},
		)

		if err != nil {
			log.Println(err)
			return err
		}

		return c.JSON(resolution)
	})

	app.Post("/resolutions/:id/complete", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return err
		}

		resolution, err := resolutionServiceClient.CompleteResolution(
			context.Background(),
			&resolutionServiceProto.CompleteResolutionRequest{
				ResolutionId: int32(id),
			},
		)

		if err != nil {
			return err
		}

		return c.JSON(resolution)
	})

	app.Post("/resolutions/:id/delete", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return err
		}

		_, err = resolutionServiceClient.DeleteResolution(
			context.Background(),
			&resolutionServiceProto.DeleteResolutionRequest{
				ResolutionId: int32(id),
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

	app.Post("/register/", func(c *fiber.Ctx) error {
		credentials := new(UserCredentails)
		if err := c.BodyParser(credentials); err != nil {
			return err
		}

		response, err := userServiceClient.CreateUser(
			context.Background(),
			&userServiceProto.UserCredentials{
				Username: credentials.Username,
				Password: credentials.Password,
			},
		)
		if err != nil {
			return err
		}

		return c.JSON(response)
	})

	app.Listen(":3000")
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
