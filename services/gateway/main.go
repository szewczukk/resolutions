package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt"
	resolutionServiceProto "github.com/szewczukk/resolution-service/proto"
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

var sampleSecretKey = []byte("SecretYouShouldHide")

func main() {
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

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/users/", func(c *fiber.Ctx) error {
		users, err := userServiceClient.GetAllUsers(
			context.Background(),
			&userServiceProto.GetAllUsersRequest{},
		)

		if err != nil {
			return err
		}

		return c.JSON(users.Users)
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

	app.Listen(":3002")
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
