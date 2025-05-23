package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	Cartid string  `gorm:"primary key" json:"cartid"`
	Cid    string  `json:"cid" validate:"required"`
	Total  float64 `json:"total"`
}

type Product_Cart struct {
	Product_cartid string `gorm:"primary key" json:"product_cartid"`
	Cartid         string `json:"cartid" validate:"required"`
	Pid            string `json:"pid" validate:"required"`
	Quantity       int    `json:"quantity" validate:"required"`
}

type Cart_obj struct {
	Cartid   string    `json:"cartid"`
	Cid      string    `json:"cid" validate:"required"`
	Products []Product `json:"products" validate:"required"`
	Total    float32   `json:"total"`
}
type GetCart struct {
	Cartid   string  `json:"cartid"`
	Cid      string  `json:"cid"`
	Pid      string  `json:"pid"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Total    float32 `json:"total"`
}

type ShowCart struct {
	Cartid   string  `json:"cartid"`
	Cid      string  `json:"cid"`
	Pid      string  `json:"pid"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Total    float32 `json:"total"`
}

type Product struct {
	Id       string  `gorm:"primary key" json:"id" `
	Name     string  `json:"name" validate:"required"`
	Desc     string  `json:"desc" validate:"required"`
	Price    float32 `json:"price" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`
}

type Customer struct {
	Cid     string `json:"cid"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Emailid string `json:"emailid" validate:"required"`
}

func MigrateCart(db *gorm.DB) error {
	err := db.AutoMigrate(&Cart{})
	err = db.AutoMigrate(&Product_Cart{})
	return err
}
