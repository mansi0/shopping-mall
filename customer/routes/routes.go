package routes

import (
	"fmt"
	"shopping-mall-customer/controller"
)

func StartRouter() {

	app := SetUpRouter()
	customer := app.Group("/customer")
	customer.Post("/createcustomer", (controller.CreateCustomer))
	customer.Get("/getcustomer", (controller.GetCustomers))
	customer.Get("/getcustomerbyid/:id", (controller.GetCustomerById))
	customer.Patch("/updatecustomer", controller.UpdateCustomerEmail)
	customer.Delete("/deletecustomer/:id", controller.DeleteCustomerById)
	//starting product router on 8082

	app.Listen(":8082")
	fmt.Println("server is listening on 8082")
}
