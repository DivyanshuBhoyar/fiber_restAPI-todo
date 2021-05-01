package controllers

import (
	"os" // new

	"github.com/DivyanshuBhoyar/fiber_ecommerce/config" // new
	"github.com/DivyanshuBhoyar/fiber_ecommerce/models"

	// new
	// new
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" // new
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// new
	// new
)

func Get_Products(c *fiber.Ctx) error { //receives a fiber context
	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION")) //specified collecton
	query := bson.D{{}}

	cursor, err := productCollection.Find(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "something went wrong",
			"error":   err.Error(),
		})
	}
	var products = make([]models.Product, 0)

	err = cursor.All(c.Context(), &products)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"err":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"products": products,
		},
	})
}

func Get_Product(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION")) //specified collecton

	paramId := c.Params("id")

	id, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse Id",
		})
	}

	product := &models.Product{}
	query := bson.D{{Key: "_id", Value: id}}

	err = productCollection.FindOne(c.Context(), query).Decode(product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"err":     err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"product": product,
		},
	})
}

func Add_Product(c *fiber.Ctx) error { //receives a fiber context
	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION")) //specified collecton

	data := new(models.Product)

	err := c.BodyParser(&data) //parse req body as per model
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"err":     err,
		})
	}
	data.ID = nil                                                 //mongodb does that automatically
	result, err := productCollection.InsertOne(c.Context(), data) //pass new instance of model to insertOne
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert todo",
			"error":   err,
		})
	}
	//to get recently inserted doc :
	product := &models.Product{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}       //result comes from inserOne
	productCollection.FindOne(c.Context(), query).Decode(product) //decode incoming as defined in product schema

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"product": product,
		},
	})
}

func Delete_Product(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION")) //specified collecton

	var paramId = c.Params("id")
	id, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
		})
	}
	query := bson.D{{Key: "_id", Value: id}}

	err = productCollection.FindOneAndDelete(c.Context(), query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusFound).JSON(fiber.Map{
				"success": false,
				"message": "Product not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete product",
			"error":   err,
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func Update_Product(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION")) //specified collecton
	product_id := c.Params("id")

	id, err := primitive.ObjectIDFromHex(product_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}

	data := new(models.Product)
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse ID",
		})
	}
	query := bson.D{{Key: "_id", Value: id}}

	//update data
	var dataToUpdate bson.D

	if data.Name != nil {
		//product.Name = *data.Name
		dataToUpdate = append(dataToUpdate, bson.E{Key: "name", Value: data.Name})
	}
	if data.Category != nil {
		//product.Category = *data.Category
		dataToUpdate = append(dataToUpdate, bson.E{Key: "category", Value: data.Category})
	}
	if data.Price != nil {
		//product.Price = *data.Price
		dataToUpdate = append(dataToUpdate, bson.E{Key: "price", Value: data.Price})
	}
	if data.Image != nil {
		//product.Image = *data.Image
		dataToUpdate = append(dataToUpdate, bson.E{Key: "image", Value: data.Image})
	}
	if data.Description != nil {
		//product.Description = *data.Description
		dataToUpdate = append(dataToUpdate, bson.E{Key: "description", Value: data.Description})
	}

	update := bson.D{
		{
			Key: "$set", Value: dataToUpdate,
		},
	}
	err = productCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Product Not found",
				"error":   err,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update todo",
			"error":   err,
		})
	}

	product := &models.Product{}
	productCollection.FindOne(c.Context(), query).Decode(product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"product": product,
		},
	})

}
