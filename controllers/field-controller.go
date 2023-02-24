package controllers

import (
    "context"
    "bola-oi/configs"
    "bola-oi/models"
    "bola-oi/responses"
    "net/http"
    "time"
    
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	// "log"
)

var fieldCollection *mongo.Collection = configs.GetCollection(configs.DB, "fields")

func CreateField(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var field models.Field
    defer cancel()

    //validate the request body
    if err := c.BodyParser(&field); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&field); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

    newField := models.Field{
		Id:       primitive.NewObjectID(),
        Name:     field.Name,
        Location: field.Location,
        Type:     field.Type,
        Price:    field.Price,
        Unit:     field.Unit,
    }

    
	_, err := fieldCollection.InsertOne(ctx, newField)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }
	
    return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"insertedId": newField.Id}})
}

func GetFieldById(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    fieldId := c.Params("fieldId")
    var field models.Field

    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(fieldId)
    err := fieldCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&field)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": field}})
}

func GetAllFields(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var fields []models.Field
    defer cancel()

    results, err := fieldCollection.Find(ctx, bson.M{})

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleField models.Field
        if err = results.Decode(&singleField); err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }

        fields = append(fields, singleField)
    }

    return c.Status(http.StatusOK).JSON(
        responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": fields}},
    )
}

func EditField(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    fieldId := c.Params("fieldId")
    var field models.Field
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(fieldId)

    //validate the request body
    if err := c.BodyParser(&field); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&field); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }

    update := bson.M{"name": field.Name, "location": field.Location, "type": field.Type, "price": field.Price, "unit": field.Unit}

    result, err := fieldCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //get updated user details
    var updatedField models.Field
    if result.MatchedCount == 1 {
        err := fieldCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedField)
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }
    }

    return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedField}})
}

func DeleteField(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    fieldId := c.Params("fieldId")
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(fieldId)

    result, err := fieldCollection.DeleteOne(ctx, bson.M{"id": objId})
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    if result.DeletedCount < 1 {
        return c.Status(http.StatusNotFound).JSON(
            responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Field with specified ID not found!"}},
        )
    }

    return c.Status(http.StatusOK).JSON(
        responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Field successfully deleted!"}},
    )
}