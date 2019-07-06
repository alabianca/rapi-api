package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/alabianca/rapi-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserCollection = "users"

type UserDAL interface {
	CreateUser(u *User) map[string]interface{}
	Login(email, password string) map[string]interface{}
	GetUserById(id primitive.ObjectID) map[string]interface{}
	AddRecord(userId primitive.ObjectID, id primitive.ObjectID) map[string]interface{}
	GetRecords(userId primitive.ObjectID) map[string]interface{}
}

type User struct {
	FirstName string               `json:"firstName"`
	LastName  string               `json:"lastName"`
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
	Records   []primitive.ObjectID `json:"records" bson:"records,omitempty"`
}

type UserSource struct{}

// Validate validates if a user by the u.Email already exists
func (us UserSource) Validate(u *User) (map[string]interface{}, bool) {
	if !strings.Contains(u.Email, "@") {
		return utils.Message(http.StatusBadRequest, "Email address is required"), false
	}

	if len(u.Password) < 6 {
		return utils.Message(http.StatusBadRequest, "Passowrd must be minimum 6 characters long"), false
	}

	// Email must be unique
	db, err := GetDB()
	filter := bson.D{{"email", u.Email}}

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not connect to database"), false
	}

	users := db.Collection(UserCollection)

	if err := users.FindOne(context.TODO(), filter).Decode(u); err != mongo.ErrNoDocuments {

		return utils.Message(http.StatusConflict, "Email already exists"), false
	}

	return utils.Message(http.StatusOK, "Requirement Passed"), true
}

// Create inserts a new user if the user with `email` does not yet exist
func (us UserSource) CreateUser(u *User) map[string]interface{} {
	if resp, ok := us.Validate(u); !ok {
		return resp
	}

	// hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on database")
	}

	users := db.Collection(UserCollection)

	insertResult, err := users.InsertOne(context.TODO(), u)

	if err != nil {
		return utils.Message(http.StatusNotModified, "Error Inserting user document.\n"+err.Error())
	}

	u.Password = "" // clear the password from ever leaving the api
	u.ID = insertResult.InsertedID.(primitive.ObjectID)

	response := utils.Message(http.StatusCreated, "Successfully created user")
	response["data"] = u

	return response
}

func (us UserSource) Login(email, password string) map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on database")
	}

	user := &User{}
	filter := bson.D{{"email", email}}
	users := db.Collection(UserCollection)

	if err := users.FindOne(context.TODO(), filter).Decode(user); err == mongo.ErrNoDocuments {
		return utils.Message(http.StatusNotFound, fmt.Sprintf("User %s not found", email))
	}

	// verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return utils.Message(http.StatusUnauthorized, "Authentication Error")
	}

	// user is legit. send up a token

	response := utils.Message(http.StatusFound, "Successfully got user")
	response["data"] = user

	return response
}

func (us UserSource) GetUserById(id primitive.ObjectID) map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on the database")
	}

	user := &User{}
	users := db.Collection(UserCollection)
	filter := bson.D{{"_id", id}}

	if err := users.FindOne(context.TODO(), filter).Decode(user); err == mongo.ErrNoDocuments {
		return utils.Message(http.StatusNotFound, "User Not Found")
	}

	// found the user!
	user.Password = ""
	response := utils.Message(http.StatusOK, "User Found")

	response["data"] = user

	return response
}

func (us UserSource) AddRecord(userId primitive.ObjectID, id primitive.ObjectID) map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on the database")
	}

	users := db.Collection(UserCollection)

	filter := bson.D{{"_id", userId}}
	update := bson.D{
		{"$push", bson.D{
			{"records", id},
		}},
	}

	_, err = users.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, fmt.Sprintf("Error Updating User %s", err))
	}

	return us.GetUserById(userId)
}

func (us UserSource) GetRecords(userId primitive.ObjectID) map[string]interface{} {
	db, err := GetDB()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on the database")
	}

	var results []*URLRecord
	records := db.Collection("records")

	filter := bson.D{{"userId", userId}}

	cur, err := records.Find(context.TODO(), filter)
	if err != nil {
		return utils.Message(http.StatusNotFound, fmt.Sprintf("No Records found for User %s", userId.Hex()))
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var record URLRecord

		if cur.Decode(&record) != nil {
			return utils.Message(http.StatusInternalServerError, fmt.Sprintf("Could Not Decode Record. Cursor Error"))
		}
		log.Printf("Appending Record %s", record.URL)
		results = append(results, &record)
	}

	if len(results) == 0 {
		return utils.Message(http.StatusNotFound, fmt.Sprintf("No Records Found for %s", userId.Hex()))
	}

	resp := utils.Message(http.StatusOK, "Records found")
	resp["data"] = results

	return resp

}
