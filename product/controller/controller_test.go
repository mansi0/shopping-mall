package controller_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"shopping-mall-product/controller"
	model "shopping-mall-product/models"
	routetest "shopping-mall-product/routes"
	"shopping-mall-product/storage"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestCreateProduct(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()

	app.Post("/product/createproduct", controller.CreateProduct)
	//415 feild missing

	product := model.Product{
		Name:  "Samsung M73",
		Desc:  "Samsung",
		Price: 30000.0,
		// Quantity: 15,
	}

	jsonData, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/product/createproduct", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusUnsupportedMediaType, resp.StatusCode, "status code")

	//422 wrong data type
	newProduct := map[string]any{
		"Name":     "Samsung M73",
		"Desc":     "Samsung",
		"Price":    "30000.0",
		"Quantity": "15",
	}

	jsonData, _ = json.Marshal(newProduct)
	req = httptest.NewRequest(http.MethodPost, "/product/createproduct", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	utils.AssertEqual(t, http.StatusUnprocessableEntity, resp.StatusCode, "status code")

	//200 product created
	product_obj := model.Product{
		Name:     "Samsung M73",
		Desc:     "Samsung",
		Price:    30000.0,
		Quantity: 15,
	}

	jsonData, _ = json.Marshal(product_obj)
	req = httptest.NewRequest(http.MethodPost, "/product/createproduct", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	data, _ := ioutil.ReadAll(resp.Body)
	expected := map[string]string{
		"message": "Product has been added",
	}
	expectedBytes, _ := json.Marshal(expected)

	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}
}

func TestUpdateQuantity(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Patch("/product/updateproduct", controller.UpdateQuantity)

	//415 feild missing

	updateProduct := model.UpdateProductQuantity{
		Id: "23e",
		// Quantity: 3,
	}

	jsonData, _ := json.Marshal(updateProduct)
	req := httptest.NewRequest(http.MethodPatch, "/product/updateproduct", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusUnsupportedMediaType, resp.StatusCode, "status code")

	//422 wrong data type
	newProduct := map[string]any{
		"Id":       "23e",
		"Quantity": "3",
	}

	jsonData, _ = json.Marshal(newProduct)

	resp, _ = app.Test(httptest.NewRequest(http.MethodPatch, "/product/updateproduct", bytes.NewBuffer(jsonData)))

	utils.AssertEqual(t, http.StatusUnprocessableEntity, resp.StatusCode, "status code")

	// 404 product does not exist
	product := model.UpdateProductQuantity{
		Id:       "23e",
		Quantity: 3,
	}

	jsonData, _ = json.Marshal(product)
	req = httptest.NewRequest(http.MethodPatch, "/product/updateproduct", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "product not found",
	}
	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	// 200 product updated
	product_obj := model.UpdateProductQuantity{
		Id:       "124",
		Quantity: 3,
	}

	jsonData, _ = json.Marshal(product_obj)
	req = httptest.NewRequest(http.MethodPatch, "/product/updateproduct", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ = ioutil.ReadAll(resp.Body)

	expected = map[string]string{
		"message": "product updated successfully",
	}
	expectedBytes, _ = json.Marshal(expected)
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

}

func TestGetProductById(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Get("/product/getproductbyid/:id", controller.GetProductById)
	//404 not found

	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/product/getproductbyid/a", nil))
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "could not get the product",
	}

	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//200 product fetched
	resp, _ = app.Test(httptest.NewRequest(http.MethodGet, "/product/getproductbyid/124", nil))
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")

}

func TestDeleteProductById(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Delete("/product/deleteproduct/:id", controller.DeleteProductById)
	//404 not found

	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/product/deleteproduct/a", nil))
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "Product not found",
	}

	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//200 product deleted succesfully
	resp, _ = app.Test(httptest.NewRequest(http.MethodDelete, "/product/deleteproduct/edc", nil))
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	data, _ = ioutil.ReadAll(resp.Body)
	expected = map[string]string{
		"message": "product deleted successfully",
	}
	expectedBytes, _ = json.Marshal(expected)
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}
}
