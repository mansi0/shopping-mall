package controller

import (
	"log"
	"net/http"
	"shopping-mall-product/common"
	"shopping-mall-product/models"
	"shopping-mall-product/storage"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/validator.v9"
)

func GetProducts(context *fiber.Ctx) error {
	productModels := []models.Product{}
	r := storage.FetchRepo()
	err := r.DB.Find(&productModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get products"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "products fetched successfully",
			"data":    productModels,
		})
	return nil

}

func GetProductById(context *fiber.Ctx) error {
	id := context.Params("id")
	product := &models.Product{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "ID cannot be empty"})
		return nil
	}
	r := storage.FetchRepo()
	err := r.DB.Where("id =?", id).First(product).Error
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "could not get the product"})
	}
	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Product id fetched successfully,", "data": product})

}

func CreateProduct(context *fiber.Ctx) error {
	var product models.Product
	randomId := common.GenerateRandomId(3)
	err := context.BodyParser(&product)
	if err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err.Error()})

	}
	//check required feilds
	v := validator.New()
	err = v.Struct(product)
	if err != nil {
		return context.Status(http.StatusUnsupportedMediaType).JSON(
			&fiber.Map{"message": err.Error()})
		// return err

	}
	r := storage.FetchRepo()
	product.Id = randomId
	err = r.DB.Create(&product).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create product"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Product has been added"})
	return nil
}

func UpdateQuantity(context *fiber.Ctx) error {
	updateProduct := models.UpdateProductQuantity{}
	fetchProduct := models.Product{}
	err := context.BodyParser(&updateProduct)
	if err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err.Error()})

	}

	//check required feilds
	v := validator.New()
	err = v.Struct(updateProduct)
	if err != nil {
		return context.Status(http.StatusUnsupportedMediaType).JSON(
			&fiber.Map{"message": err.Error()})

	}

	r := storage.FetchRepo()
	//check product id exist or not
	err = r.DB.Where("id=?", updateProduct.Id).First(&fetchProduct).Error
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "product not found",
		})

	}
	err = r.DB.Model(&models.Product{}).Where("id=?", updateProduct.Id).Update("quantity", updateProduct.Quantity).Error
	if err != nil {
		log.Println(err)
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update the product"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "product updated successfully"})
	return nil
}

func DeleteProductById(context *fiber.Ctx) error {
	product := models.Product{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	r := storage.FetchRepo()
	//check product id exist or not
	err := r.DB.Where("id=?", id).First(&product).Error
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "Product not found",
		})

	}
	err = r.DB.Where("id =?", id).Delete(&product).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "Product not found"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "product deleted successfully"})
	return nil
}
