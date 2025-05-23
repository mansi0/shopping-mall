package controller

import (
	"fmt"
	"log"
	"net/http"
	"shopping-mall-customer/common"
	models "shopping-mall-customer/models"
	"shopping-mall-customer/storage"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/validator.v9"
)

func GetCustomers(context *fiber.Ctx) error {
	customers := []models.Customer{}
	r := storage.FetchRepo()
	err := r.DB.Find(&customers).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get customers"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Customers fetched successfully",
			"data":    customers,
		})
	return nil

}

func GetCustomerById(context *fiber.Ctx) error {
	id := context.Params("id")
	log.Print(id)
	fmt.Printf("id issssssss: %v\n", id)
	customer := &models.Customer{}
	if id == "" {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "ID cannot be empty"})

	}
	r := storage.FetchRepo()
	err := r.DB.Where("cid =?", id).First(customer).Error
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "could not get the customer"})

	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Customer details fetched successfully,",
		"data":    customer,
	})
	return nil

}

func CreateCustomer(context *fiber.Ctx) error {
	var customer models.Customer
	err := context.BodyParser(&customer)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "invalid input"})
		return err
	}

	//check required feilds
	v := validator.New()
	err = v.Struct(customer)
	if err != nil {
		return context.Status(http.StatusUnsupportedMediaType).JSON(
			&fiber.Map{"message": err})

	}

	randomId := common.GenerateRandomId(3)
	customer.Cid = randomId
	r := storage.FetchRepo()
	err = r.DB.Create(&customer).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create customer object"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Customer has been added"})
	return nil
}

func UpdateCustomerEmail(context *fiber.Ctx) error {
	updateCustomer := models.UpdateCustomerEmail{}
	fetchCustomer := models.Customer{}
	err := context.BodyParser(&updateCustomer)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "invalid input"})
		return err
	}

	//check required feilds
	v := validator.New()
	err = v.Struct(updateCustomer)
	if err != nil {
		return context.Status(http.StatusUnsupportedMediaType).JSON(
			&fiber.Map{"message": err})

	}

	r := storage.FetchRepo()
	//check customer id exist or not
	err = r.DB.Where("cid=?", updateCustomer.Cid).First(&fetchCustomer).Error
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "customer not found",
		})

	}

	//update customer
	err = r.DB.Model(&models.Customer{}).Where("cid=?", updateCustomer.Cid).Update("emailid", updateCustomer.EmailId).Error
	if err != nil {
		log.Println(err)
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update the customer"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "customer updated successfully"})
	return nil
}

func DeleteCustomerById(context *fiber.Ctx) error {
	customer := models.Customer{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	r := storage.FetchRepo()

	//check id exist or not
	err := r.DB.Where("cid=?", id).First(&customer).Error
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "Customer not found",
		})

	}
	err = r.DB.Where("cid =?", id).Delete(&models.Customer{}).Error
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete customer"})

	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Customer deleted successfully"})
	return nil
}
