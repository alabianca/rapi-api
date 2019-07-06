package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alabianca/rapi-api/models"
	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateUserHandler(t *testing.T) {
	// get instance of test api
	api := getTestAPI()
	// create the body of the request
	user := &models.Registration{
		FirstName: "Tester",
		LastName:  "Le Tester",
		Email:     "test@letester.com",
		Password:  "1234",
		Verify:    "1234",
	}
	body, err := json.Marshal(user)

	if err != nil {
		t.Fatal(err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	// create a response recorder. which satisfies the ResponseWriter interface
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(api.CreateUser)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Expected %d, but got %d\n", http.StatusCreated, status)
	}

}

func TestCreateUserHandlerWithNoMatchPassword(t *testing.T) {
	// get instance of test api
	api := getTestAPI()
	// create the body of the request
	user := &models.Registration{
		FirstName: "Tester",
		LastName:  "Le Tester",
		Email:     "test@letester.com",
		Password:  "1234",
		Verify:    "12345",
	}
	body, err := json.Marshal(user)

	if err != nil {
		t.Fatal(err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
	// create a response recorder. which satisfies the ResponseWriter interface
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(api.CreateUser)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected %d, but got %d", http.StatusBadRequest, status)
	}
}

type TestUsersDAL struct {
}

func (t TestUsersDAL) CreateUser(u *models.User) map[string]interface{} {
	response := utils.Message(http.StatusCreated, "User Created")
	u.Password = ""
	response["data"] = u

	return response
}
func (t TestUsersDAL) Login(email, password string) map[string]interface{} {
	response := utils.Message(http.StatusOK, "Login Success")
	response["data"] = &models.User{
		Email: email,
	}

	return response
}

func (t TestUsersDAL) GetUserById(id primitive.ObjectID) map[string]interface{} {
	response := utils.Message(http.StatusOK, "OK")
	response["data"] = &models.User{
		ID:    id,
		Email: "Tester.test@test.com",
	}

	return response
}

func (t TestUsersDAL) AddRecord(userId primitive.ObjectID, id primitive.ObjectID) map[string]interface{} {
	return nil
}
func (t TestUsersDAL) GetRecords(userId primitive.ObjectID) map[string]interface{} {
	return nil
}
