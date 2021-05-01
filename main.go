package main

import (
	"log"

	"github.com/DivyanshuBhoyar/fiber_ecommerce/config"
	"github.com/DivyanshuBhoyar/fiber_ecommerce/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv" // new
)

func setup_routes(app *fiber.App) { //receives app instance as it is called directly from main func

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Server live 📡",
			"root":    true,
		})
	})

	api := app.Group("/api") //group
	routes.ProductRoute(api) //router instance (graup ) is passed
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())

	setup_routes(app) //instance of app shared with function setup routes
	err = app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
