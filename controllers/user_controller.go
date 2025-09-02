package controllers

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"

    "github.com/Ashu-300/golang-fiber-mvc-starter/database"
    "github.com/Ashu-300/golang-fiber-mvc-starter/models"
)

const SecretKey = "secret"

func Login(c *fiber.Ctx) error {
    var data map[string]string
    if err := c.BodyParser(&data); err != nil {
        return err
    }

    var user models.User
    collection := database.DB.Collection("users")
    err := collection.FindOne(context.Background(), bson.M{"email": data["email"]}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User not found",
        })
    }

    if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Invalid credentials",
        })
    }

    claims := jwt.MapClaims{
        "id":  user.ID.Hex(),
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(SecretKey))
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

    c.Cookie(&fiber.Cookie{
        Name:     "jwt",
        Value:    t,
        Expires:  time.Now().Add(time.Hour * 24),
        HTTPOnly: true,
    })

    return c.JSON(fiber.Map{"message": "success"})
}

func GetUser(c *fiber.Ctx) error {
    cookie := c.Cookies("jwt")
    token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
        return []byte(SecretKey), nil
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "unauthenticated",
        })
    }

    claims := token.Claims.(jwt.MapClaims)
    id, _ := primitive.ObjectIDFromHex(claims["id"].(string))

    var user models.User
    collection := database.DB.Collection("users")
    err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User not found",
        })
    }

    return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
    c.Cookie(&fiber.Cookie{
        Name:     "jwt",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HTTPOnly: true,
    })

    return c.JSON(fiber.Map{"message": "logged out"})
}
