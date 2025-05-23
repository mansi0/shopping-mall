package main

import (
	"shopping-mall-cart/routes"
	"shopping-mall-cart/storage"
)

func main() {
	storage.FetchRepo()
	routes.StartRouter()
}
