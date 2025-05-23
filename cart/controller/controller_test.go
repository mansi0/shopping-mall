package controller_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"shopping-mall-cart/controller"
	"shopping-mall-cart/models"
	model "shopping-mall-cart/models"
	routetest "shopping-mall-cart/routes"
	"shopping-mall-cart/storage"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestCreateCart(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()

	app.Post("/cart/createcart", controller.CreateCart)
	//415 feild missing
	cart := model.Cart_obj{
		// Cid: "23e", cid missing
		Products: []models.Product{
			{
				Id:       "c44",
				Quantity: 2,
			},
			{
				Id:       "bf5",
				Quantity: 3,
			},
		},
	}

	jsonData, _ := json.Marshal(cart)
	req := httptest.NewRequest(http.MethodPost, "/cart/createcart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusUnsupportedMediaType, resp.StatusCode, "status code")

	//404 check product exist or not
	cartobj := model.Cart_obj{
		Cid: "23e",
		Products: []models.Product{
			{
				Id:       "c44",
				Quantity: 2,
			},
			{
				Id:       "bf5",
				Quantity: 3,
			},
		},
	}

	jsonData, _ = json.Marshal(cartobj)
	req = httptest.NewRequest(http.MethodPost, "/cart/createcart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "invalid input - product not exist",
	}
	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//404 check customer exist or not

	cart_obj := model.Cart_obj{
		Cid: "23e",
		Products: []models.Product{
			{
				Id:       "124",
				Quantity: 1,
			},
			{
				Id:       "9d0",
				Quantity: 1,
			},
		},
	}

	jsonData, _ = json.Marshal(cart_obj)
	req = httptest.NewRequest(http.MethodPost, "/cart/createcart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ = ioutil.ReadAll(resp.Body)

	expected = map[string]string{
		"message": "invalid input- customer not exist",
	}

	expectedBytes, _ = json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")

	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//404 check product in stock or not

	cartobj2 := model.Cart_obj{
		Cid: "e10",
		Products: []models.Product{
			{
				Id:       "124",
				Quantity: 20,
			},
			{
				Id:       "9d0",
				Quantity: 30,
			},
		},
	}

	jsonData, _ = json.Marshal(cartobj2)
	req = httptest.NewRequest(http.MethodPost, "/cart/createcart", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ = ioutil.ReadAll(resp.Body)

	expected = map[string]string{
		"message": "Product out of stock",
	}

	expectedBytes, _ = json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")

	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	// 200 cart created
	/*
		cartobj3 := model.Cart_obj{
			Cid: "e10",
			Products: []models.Product{
				{
					Id:       "6c5",
					Quantity: 1,
				},
				{
					Id:       "7b4",
					Quantity: 1,
				},
			},
		}

		jsonData, _ := json.Marshal(cartobj3)
		req := httptest.NewRequest(http.MethodPost, "/cart/createcart", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
		data, _ := ioutil.ReadAll(resp.Body)

		expected := map[string]string{
			"message": "Cart has been added",
		}

		expectedBytes, _ := json.Marshal(expected)

		if string(data) != string(expectedBytes) {
			t.Errorf("expected %v got %v", string(expectedBytes), string(data))
		}*/

}

func TestGetCartById(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Get("/cart/getcartbyid/:id", controller.GetCartById)
	//404 not found

	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/cart/getcartbyid/a", nil))
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "cart id doesnot exist",
	}

	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//200 cart fetched
	resp, _ = app.Test(httptest.NewRequest(http.MethodGet, "/cart/getcartbyid/973", nil))
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")

}

func TestDeleteCartById(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Delete("/cart/deletecart/:id", controller.DeleteCartById)
	//404 not found

	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/cart/deletecart/a", nil))
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "Cart not found",
	}

	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//200 cart deleted succesfully
	resp, _ = app.Test(httptest.NewRequest(http.MethodDelete, "/cart/deletecart/bf5", nil))
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	data, _ = ioutil.ReadAll(resp.Body)
	expected = map[string]string{
		"message": "cart deleted successfully",
	}
	expectedBytes, _ = json.Marshal(expected)
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}
}
