package model

import "gorm.io/gorm"

type Customer struct {
	Cid     string `json:"cid"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Emailid string `json:"emailid" validate:"required"`
}

type UpdateCustomerEmail struct {
	Cid     string `json:"cid" validate:"required"`
	EmailId string `json:"emailid" validate:"required"`
}

func MigrateCustomer(db *gorm.DB) error {
	err := db.AutoMigrate(&Customer{})
	return err
}
