package common

import (
	"errors"
	model "shopping-mall-customer/models"
	"shopping-mall-customer/storage"

	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomId(length int) string {
	// generating random string of 6 length
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
func IsExistCustomerId(cid string) (bool, error) {
	r := storage.FetchRepo()
	customer := model.Customer{}
	err := r.DB.First(&customer, "cid=?", cid).Error
	if err != nil {
		return false, errors.New("invalid input- customer not exist")
	}
	return true, nil
}
