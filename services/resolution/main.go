package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Resolution struct {
	gorm.Model
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

type CreateResolutionDTO struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Resolution{})

	app := fiber.New()

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

		resolution := &Resolution{
			Name:   createResolutionDto.Name,
			UserId: createResolutionDto.UserId,
		}

		db.Create(&resolution)

		var resolutions []Resolution
		db.Find(&resolutions)

		return c.JSON(resolutions)
	})

	app.Listen(":3001")
}
