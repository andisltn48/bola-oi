package controllers

import (
	"context"
    "bola-oi/configs"
    "bola-oi/models"
    "bola-oi/responses"
    "net/http"
    "time"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
	// jwtware "github.com/gofiber/jwt/v3"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()
const jwtSecret = "asecret"

func Register(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var user models.User
    defer cancel()

    //validate the request body
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

    newUser := models.User{
		Id:       	  primitive.NewObjectID(),
        Name:     	  user.Name,
        Email: 		  user.Email,
        Password:     user.Password,
    }

    
	_, err := userCollection.InsertOne(ctx, newUser)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }
	
    return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"insertedId": newUser.Id}})
}

func Login(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    type Request struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    var req Request
    var user models.User
    defer cancel()

    //validate the request body
    if err := c.BodyParser(&req); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }
    
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "email not found"}})
    }
	
    if req.Password != user.Password {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "password wrong"}})
    }

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["sub"] = user.Id
    claims["exp"] = time.Now().Add(time.Hour * 24 * 7)

    t, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data":err.Error()}})
    }
   
    return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{
        "token": t,
        "email": req.Email,
    }})
}