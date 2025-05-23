package routes

import "shopping-mall-cart/controller"

func StartRouter() {

	app := SetUpRouter()
	cart := app.Group("/cart")
	cart.Post("/createcart", (controller.CreateCart))
	cart.Get("getcart", controller.GetCart)
	cart.Get("/getcart/:id", controller.GetCartById)
	cart.Put("/updatecart", controller.UpdateCart)
	cart.Delete("/deletecart/:id", controller.DeleteCartById)
	cart.Delete("/checkout/:id", controller.Checkout)
	//starting cart router on 8080
	app.Listen(":8080")
}
