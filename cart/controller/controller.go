package controller

import (
	"errors"
	"log"
	"net/http"
	"shopping-mall-cart/common"
	customerCommon "shopping-mall-cart/common"
	productCommon "shopping-mall-cart/common"
	"shopping-mall-cart/models"
	"shopping-mall-cart/storage"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/validator.v9"
)

func GetCartById(context *fiber.Ctx) error {
	id := context.Params("id")
	var cart_obj models.Cart_obj
	GetCart := []models.GetCart{}
	fetchCart := models.Cart{}
	if id == "" {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "ID cannot be empty"})

	}
	r := storage.FetchRepo()
	//check cart id exist or not
	err := r.DB.First(&fetchCart, "cartid=?", id).Error
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "cart id doesnot exist",
		})

	}
	err = r.DB.Raw("select carts.cartid,cid,pid,name,price,product_carts.quantity,total from carts,products,product_carts where carts.cartid = product_carts.cartid and product_carts.pid = products.id and carts.cartid =?;", id).Scan(&GetCart).Error
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "could not get the cart"})

	}
	for _, cart := range GetCart {
		if cart_obj.Cartid == "" {
			cart_obj.Cartid = cart.Cartid
			cart_obj.Cid = cart.Cid
			product := models.Product{

				Name:     cart.Name,
				Price:    cart.Price,
				Quantity: cart.Quantity,
			}
			cart_obj.Products = append(cart_obj.Products, product)

		} else {
			product := models.Product{

				Name:     cart.Name,
				Price:    cart.Price,
				Quantity: cart.Quantity,
			}
			cart_obj.Products = append(cart_obj.Products, product)
		}
		cart_obj.Total = cart.Total
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Cart details fetched successfully,",
		"data":    cart_obj,
	})
	return nil

}

func GetCart(context *fiber.Ctx) error {
	//id := context.Params("id")

	GetCart := []models.GetCart{}
	fetchCart := models.Cart{}
	// if id == "" {
	// 	return context.Status(http.StatusInternalServerError).JSON(
	// 		&fiber.Map{
	// 			"message": "ID cannot be empty"})

	// }
	r := storage.FetchRepo()
	//check cart id exist or not
	err := r.DB.Find(&fetchCart).Error
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "cart id doesnot exist",
		})

	}
	err = r.DB.Raw("select carts.cartid,cid,pid,name,price,product_carts.quantity,total from carts,products,product_carts where carts.cartid = product_carts.cartid and product_carts.pid = products.id;").Scan(&GetCart).Error
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "could not get the cart"})

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Cart details fetched successfully,",
		"data":    GetCart,
	})
	return nil

}

func CreateCart(context *fiber.Ctx) error {
	var cart_obj models.Cart_obj
	var cart models.Cart
	var product_cart []models.Product_Cart
	err := context.BodyParser(&cart_obj)
	if err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err.Error()})
	}
	//check required feilds
	v := validator.New()
	err = v.Struct(cart_obj)
	if err != nil {
		return context.Status(http.StatusUnsupportedMediaType).JSON(
			&fiber.Map{"message": err.Error()})

	}

	//check product exist or not
	for _, productUnit := range cart_obj.Products {
		checkProductExist, err := productCommon.CheckProductExist(productUnit.Id)
		if err != nil {
			return context.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": err.Error()})

		} else if !checkProductExist {
			return context.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "invalid input - product not exist"})
		}
	}

	//check customer exist or not
	checkCustomerExist, err := customerCommon.IsExistCustomerId(cart_obj.Cid)
	if err != nil {
		return context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": err.Error()})

	} else if !checkCustomerExist {
		return context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "invalid input- customer not exist"})
	}
	randomId := common.GenerateRandomId(3)
	//check product in stock or not
	for _, productUnit := range cart_obj.Products {

		checkQuantity, err := productCommon.CheckQuantity(productUnit.Id, productUnit.Quantity)
		log.Println(productUnit, checkQuantity)
		if err != nil {
			context.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": err})
			return err
		} else if !checkQuantity {
			return context.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "Product out of stock"})
		} else {
			//make product obj
			product_cart_unit := models.Product_Cart{}
			product_cart_unit.Cartid = randomId
			product_cart_unit.Pid = productUnit.Id
			product_cart_unit.Quantity = productUnit.Quantity
			product_cart = append(product_cart, product_cart_unit)

			// calculate total
			price, err := productCommon.GetPrice(productUnit.Id)
			if err != nil {
				context.Status(http.StatusInternalServerError).JSON(
					&fiber.Map{"message": err})
				return err
			}
			cart.Total += float64(price * float32(productUnit.Quantity))
			//decrement count of that product
			// err = productCommon.UpdateProductQuantity(productUnit.Id, productUnit.Quantity)
		}
	}

	cart.Cartid = randomId
	cart_obj.Cartid = cart.Cartid
	cart.Cid = cart_obj.Cid

	r := storage.FetchRepo()
	err = r.DB.Create(&cart).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not create cart object"})
		return err
	}
	err = r.DB.Create(&product_cart).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not create product-cart object"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Cart has been added",
			"data":    cart_obj,
		})
	return nil
}

func UpdateCart(context *fiber.Ctx) error {
	UpdateCart := models.Product_Cart{}
	fetchCart := models.Cart{}
	err := context.BodyParser(&UpdateCart)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "invalid input"})
		return err
	}

	//check required feilds
	v := validator.New()
	err = v.Struct(UpdateCart)
	if err != nil {
		context.Status(http.StatusUnsupportedMediaType).JSON(
			&fiber.Map{"message": err})
		return err
	}

	r := storage.FetchRepo()
	//check cart id exist or not
	err = r.DB.First(&fetchCart, "cartid=?", UpdateCart.Cartid).Error
	if err != nil {
		context.Status(fiber.ErrInternalServerError.Code).JSON(&fiber.Map{
			"message": err,
		})
		return errors.New("cart id doesnot exist")
	}

	//check product Exist or not
	checkProductExist, err := productCommon.CheckProductExist(UpdateCart.Pid)
	if err != nil {
		return context.Status(fiber.ErrInternalServerError.Code).JSON(&fiber.Map{
			"message": "Product not exist",
		})

	} else if !checkProductExist {
		return context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": " Product not exist"})
	}

	//check product in stock or not

	checkQuantity, err := productCommon.CheckQuantity(UpdateCart.Pid, UpdateCart.Quantity)
	if err != nil {
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": err.Error()})
		return err
	} else if !checkQuantity {
		return context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "Product out of stock"})
	} else {
		// calculate total and update
		price, err := productCommon.GetPrice(UpdateCart.Pid)
		if err != nil {
			context.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": err})
			return err
		}
		fetchCart.Total += float64(price * float32(UpdateCart.Quantity))
		err = r.DB.Model(&fetchCart).Where("cartid=?", UpdateCart.Cartid).Update("total", fetchCart.Total).Error
		//decrement count of that product
		// err = productCommon.UpdateProductQuantity(UpdateCart.Pid, UpdateCart.Quantity)
	}

	err = r.DB.Create(&UpdateCart).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not create product-cart object"})
		return err
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Cart has been updated"})
	return nil
}

func DeleteCartById(context *fiber.Ctx) error {
	cart := models.Cart{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	r := storage.FetchRepo()

	//check id exist or not
	err := r.DB.Where("cartid=?", id).First(&cart).Error
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "Cart not found",
		})

	}
	err = r.DB.Where("cartid =?", id).Delete(&cart).Error
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete cart"})

	}
	err = r.DB.Where("cartid =?", id).Delete(&models.Product_Cart{}).Error
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete cart"})

	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "cart deleted successfully"})
	return nil
}

func Checkout(context *fiber.Ctx) error {
	//cart := models.Cart{}
	fetchCart := models.Cart{}
	productCart := []models.Product_Cart{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	r := storage.FetchRepo()

	//check cart id exist or not
	err := r.DB.First(&fetchCart, "cartid=?", id).Error
	if err != nil {
		context.Status(fiber.ErrInternalServerError.Code).JSON(&fiber.Map{
			"message": err,
		})
		return errors.New("cart id doesnot exist")
	}
	err = r.DB.Find(&productCart, "cartid=?", id).Error
	for _, product := range productCart {
		// decrement count of that product
		err = productCommon.UpdateProductQuantity(product.Pid, product.Quantity)
	}

	err = r.DB.Where("cartid =?", id).Delete(&fetchCart).Error
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete cart"})

	}
	err = r.DB.Where("cartid =?", id).Delete(&models.Product_Cart{}).Error
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete cart"})

	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "product count decremented, cart deleted successfully"})
	return nil
}
