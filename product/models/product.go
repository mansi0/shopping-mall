package models

import "gorm.io/gorm"

type Product struct {
	Id       string  `gorm:"primary key" json:"id" `
	Name     string  `json:"name" validate:"required"`
	Desc     string  `json:"desc" validate:"required"`
	Price    float32 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`
}

type UpdateProductQuantity struct {
	Id       string `json:"id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

func MigrateProduct(db *gorm.DB) error {
	err := db.AutoMigrate(&Product{})
	return err
}
