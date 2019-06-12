package models

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alabianca/rapi-api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/atrox/haikunatorgo"
)

type Personal struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Objective string `json:"objective"`
}

type Education struct {
	School    string    `json:"school"`
	Degree    string    `json:"degree"`
	GPA       float32   `json:"gpa"`
	StartDate time.Time `json:"fromDate"`
	EndDate   time.Time `json:"toDate"`
}

type Experience struct {
	Company         string    `json:"company"`
	Title           string    `json:"title"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	CurrentJob      bool      `json:"current"`
	Accomplishments []string  `json:"accomplishments"`
}

type Project struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Link   string `json:"link"`
}

type Resume struct {
	Personal    Personal
	Education   Education
	Experiences []Experience
	Projects    []Project
	Skills      []string
}

type URLRecord struct {
	ID  primitive.ObjectID `json:"id",bson:"_id,omitempty"`
	URL string             `json:"url"`
}

func (u *URLRecord) Create() map[string]interface{} {
	if u.URL == "" {
		return utils.Message(http.StatusBadRequest, "Illegal URL")
	}

	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could Not Get Hanlde On DB")
	}

	records := db.Collection("records")

	res, err := records.InsertOne(context.TODO(), u)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Error inserting the record")
	}

	u.ID = res.InsertedID.(primitive.ObjectID) // update the id for the payload

	response := utils.Message(http.StatusCreated, "Record Inserted")
	response["data"] = u

	return response
}

func GenerateRandomURL() map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could Not Get Handle On DB")
	}

	records := db.Collection("records")
	record := &URLRecord{}

	if err := generateAttempt(records, 0, 3, record); err != nil {
		return utils.Message(http.StatusConflict, "Could Not Generate a URL")
	}

	return record.Create()
}

func generateAttempt(urlCollection *mongo.Collection, current, max int, record *URLRecord) error {
	if current == max {
		return fmt.Errorf("Could not generate a random url in %d tries", max)
	}
	record.URL = makeURL()

	filter := bson.D{{"url", record.URL}}
	if err := urlCollection.FindOne(context.TODO(), filter).Decode(record); err == mongo.ErrNoDocuments {
		return nil
	}

	current++
	return generateAttempt(urlCollection, current, max, record)
}

func makeURL() string {
	host := os.Getenv("HOST")
	base := haikunator.New().Haikunate()

	return base + "." + host
}
