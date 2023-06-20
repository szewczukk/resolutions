package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	resolutionServiceProto "github.com/szewczukk/resolution-service/proto"
	userServiceProto "github.com/szewczukk/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateServicePayload struct {
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

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

	app.Listen(":3002")
}
