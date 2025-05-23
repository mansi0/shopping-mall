package common

import (
	"errors"
	"fmt"
	"math/rand"
	"shopping-mall-cart/models"
	"shopping-mall-cart/storage"
	"time"
)

func GenerateRandomId(length int) string {
	// generating random string of 6 length
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

//product

// check product exist
func CheckProductExist(id string) (bool, error) {
	r := storage.FetchRepo()
	product := models.Product{}
	err := r.DB.First(&product, "id=?", id).Error
	if err != nil {
		return false, errors.New("invalid input - product not exist")
	}
	return true, nil
}

// check product in stock or not
func CheckQuantity(id string, n int) (bool, error) {
	r := storage.FetchRepo()
	product := models.Product{}
	err := r.DB.First(&product, "id=?", id).Error
	if err != nil {
		return false, errors.New("Product out of stock")
	} else if product.Quantity-n > 0 {
		return true, nil
	} else if product.Quantity < n {
		errStr := fmt.Sprintf("%v with name : %v has stock of only %v", product.Id, product.Name, product.Quantity)
		return false, errors.New(errStr)
	}
	return false, nil
}

// get price from productId
func GetPrice(id string) (float32, error) {
	r := storage.FetchRepo()
	product := models.Product{}
	err := r.DB.Find(&product).Where("id=?", id).Error
	if err != nil {
		return 0, err
	}
	return product.Price, nil
}

// update product quantity
func UpdateProductQuantity(id string, quant int) error {
	r := storage.FetchRepo()
	product := models.Product{}
	err := r.DB.Find(&product).Where("id=?", id).Error
	err = r.DB.Model(&product).Where("id=?", id).Update("quantity", product.Quantity-quant).Error
	if err != nil {
		return err
	}
	return nil
}

//customer

func IsExistCustomerId(cid string) (bool, error) {
	r := storage.FetchRepo()
	customer := models.Customer{}
	err := r.DB.First(&customer, "cid=?", cid).Error
	if err != nil {
		return false, errors.New("invalid input- customer not exist")
	}
	return true, nil
}
