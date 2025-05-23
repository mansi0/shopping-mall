package controller_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"shopping-mall-customer/controller"
	model "shopping-mall-customer/models"
	routetest "shopping-mall-customer/routes"
	"shopping-mall-customer/storage"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestCreateCustomer(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()

	//415 feild missing
	app.Post("/customer/createcustomer", controller.CreateCustomer)

	customer := model.Customer{
		Name:    "abc",
		Address: "xyz",
		// Emailid: "abc@gmail.com",
	}

	jsonData, _ := json.Marshal(customer)
	req := httptest.NewRequest(http.MethodPost, "/customer/createcustomer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	log.Print("3")

	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusUnsupportedMediaType, resp.StatusCode, "status code")

	//422 wrong data type
	newCustomer := map[string]any{
		"Name":    "abc",
		"Address": "xyz",
		"Emailid": 1,
	}

	jsonData, _ = json.Marshal(newCustomer)
	req = httptest.NewRequest(http.MethodPost, "/customer/createcustomer", bytes.NewBuffer(jsonData))
	// req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	utils.AssertEqual(t, http.StatusUnprocessableEntity, resp.StatusCode, "status code")

	//200 customer created
	customer_obj := model.Customer{
		Name:    "mansi",
		Address: "www",
		Emailid: "mansi@gmail.com",
	}

	jsonData, _ = json.Marshal(customer_obj)
	req = httptest.NewRequest(http.MethodPost, "/customer/createcustomer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	data, _ := ioutil.ReadAll(resp.Body)
	expected := map[string]string{
		"message": "Customer has been added",
	}
	expectedBytes, _ := json.Marshal(expected)

	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}
}

func TestUpdateCustomerEmail(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Patch("/customer/updatecustomer", controller.UpdateCustomerEmail)

	//415 feild missing

	updateCustomer := model.UpdateCustomerEmail{
		Cid: "23e",
		//Emailid: "abc@gmail.com",
	}

	jsonData, _ := json.Marshal(updateCustomer)
	req := httptest.NewRequest(http.MethodPatch, "/customer/updatecustomer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusUnsupportedMediaType, resp.StatusCode, "status code")

	//422 wrong data type
	newCustomer := map[string]any{
		"Cid":     "c12",
		"Emailid": 1,
	}

	jsonData, _ = json.Marshal(newCustomer)

	resp, _ = app.Test(httptest.NewRequest(http.MethodPatch, "/customer/updatecustomer", bytes.NewBuffer(jsonData)))

	utils.AssertEqual(t, http.StatusUnprocessableEntity, resp.StatusCode, "status code")

	// 404 customer does not exist
	customer := model.UpdateCustomerEmail{
		Cid:     "cer",
		EmailId: "abc@gmail.com",
	}

	jsonData, _ = json.Marshal(customer)
	req = httptest.NewRequest(http.MethodPatch, "/customer/updatecustomer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "customer not found",
	}
	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	// 200 customer updated
	customer_obj := model.UpdateCustomerEmail{
		Cid:     "ce8",
		EmailId: "abc@gmail.com",
	}

	jsonData, _ = json.Marshal(customer_obj)
	req = httptest.NewRequest(http.MethodPatch, "/customer/updatecustomer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	data, _ = ioutil.ReadAll(resp.Body)

	expected = map[string]string{
		"message": "customer updated successfully",
	}
	expectedBytes, _ = json.Marshal(expected)
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

}
func TestGetCustomerById(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Get("/customer/getcustomerbyid/:id", controller.GetCustomerById)
	//404 not found

	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/customer/getcustomerbyid/a", nil))
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "could not get the customer",
	}

	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//200 customer fetched
	resp, _ = app.Test(httptest.NewRequest(http.MethodGet, "/customer/getcustomerbyid/77d", nil))
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")

}

func TestDeleteCustomerById(t *testing.T) {
	app := routetest.SetUpRouter()
	storage.InitDbForTesting()
	app.Delete("/customer/deletecustomer/:id", controller.DeleteCustomerById)
	//404 not found

	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/customer/deletecustomer/a", nil))
	data, _ := ioutil.ReadAll(resp.Body)

	expected := map[string]string{
		"message": "Customer not found",
	}

	expectedBytes, _ := json.Marshal(expected)
	utils.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "status code")
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}

	//200 customer deleted succesfully
	resp, _ = app.Test(httptest.NewRequest(http.MethodDelete, "/customer/deletecustomer/2ff", nil))
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	data, _ = ioutil.ReadAll(resp.Body)
	expected = map[string]string{
		"message": "Customer deleted successfully",
	}
	expectedBytes, _ = json.Marshal(expected)
	if string(data) != string(expectedBytes) {
		t.Errorf("expected %v got %v", string(expectedBytes), string(data))
	}
}
