package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Resolution struct {
	gorm.Model
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

type CreateResolutionDTO struct {
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Resolution{})

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		var users []Resolution
		db.Find(&users)

		return c.JSON(users)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		createResolutionDto := new(CreateResolutionDTO)

		if err := c.BodyParser(createResolutionDto); err != nil {
			return err
		}

		response, err := http.Get(fmt.Sprintf("http://localhost:3000/%v/", createResolutionDto.UserId))
		if err != nil {
			return err
		}

		if response.StatusCode != 200 {
			return errors.New("user not found")
		}

		resolution := &Resolution{
			Name:   createResolutionDto.Name,
			UserId: createResolutionDto.UserId,
		}

		db.Create(&resolution)

		return c.JSON(resolution)
	})

	app.Listen(":3001")
}
