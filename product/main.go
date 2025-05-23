package main

import (
	"shopping-mall-product/routes"
	"shopping-mall-product/storage"
)

func main() {
	storage.FetchRepo()
	routes.StartRouter()
}
