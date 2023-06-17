package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		var users []User
		db.Find(&users)

		return c.JSON(users)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		createUserDto := new(CreateUserDTO)

		if err := c.BodyParser(createUserDto); err != nil {
			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(createUserDto.Password), bcrypt.DefaultCost,
		)

		if err != nil {
			return err
		}

		user := &User{
			Username: createUserDto.Username,
			Password: string(hashedPassword),
		}

		db.Create(&user)

		var users []User
		db.Find(&users)

		return c.JSON(users)
	})

	app.Get("/:userId/", func(c *fiber.Ctx) error {
		userId := c.Params("userId")

		var user User
		result := db.First(&user, userId)

		if result.Error != nil {
			return result.Error
		}

		return c.JSON(user)
	})

	app.Listen(":3000")
}
