package main

import (
	"shopping-mall-customer/routes"
	"shopping-mall-customer/storage"
)

func main() {
	storage.FetchRepo()

	routes.StartRouter()

}
