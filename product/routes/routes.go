package routes

import "shopping-mall-product/controller"

func StartRouter() {

	app := SetUpRouter()
	product := app.Group("/product")
	product.Post("/createproduct", (controller.CreateProduct))
	product.Get("/getproduct", (controller.GetProducts))
	product.Get("/getproduct/:id", (controller.GetProductById))
	product.Patch("/updateproduct", controller.UpdateQuantity)
	product.Delete("/deleteproduct/:id", controller.DeleteProductById)
	//starting product router on 8081
	app.Listen(":8081")
}
