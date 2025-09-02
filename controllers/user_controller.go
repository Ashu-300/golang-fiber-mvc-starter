package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ashu-300/golang-fiber-mvc-starter/database"
	"github.com/Ashu-300/golang-fiber-mvc-starter/models"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		ID:       primitive.NewObjectID(),
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	collection := database.DB.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create user",
		})
	}

	return c.JSON(user)
}

//show more
