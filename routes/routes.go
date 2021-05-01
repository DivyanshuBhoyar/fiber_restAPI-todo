package routes

import (
	"github.com/DivyanshuBhoyar/fiber_ecommerce/controllers"
	"github.com/gofiber/fiber/v2"
)

func ProductRoute(route fiber.Router) { //receives Roouter instandce
	route.Get("/products", controllers.Get_Products)
	route.Get("/products/:id", controllers.Get_Product)
	route.Post("/products", controllers.Add_Product)
	// route.Delete("/products/:id", controllers.Delete_Product)
	route.Patch("/products/:id", controllers.Update_Product)

}
