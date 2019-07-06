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

	haikunator "github.com/atrox/haikunatorgo"
)

const ResumeCollection = "resume"

type ResumeDAL interface {
	CreateResume(r *Resume) map[string]interface{}
	GetResumes(userId primitive.ObjectID) map[string]interface{}
	GetResumeByID(id primitive.ObjectID) map[string]interface{}
}

type Resume struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId"`
	Personal    Personal           `json:"personal"`
	Education   Education          `json:"education"`
	Experiences []Experience       `json:"experience"`
	Projects    []Project          `json:"projects"`
	Skills      []string           `json:"skills"`
}

func CreateResume(r *Resume) map[string]interface{} {
	db, err := GetDB()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could Not Connect to DB")
	}

	resume := db.Collection(ResumeCollection)

	insertResult, err := resume.InsertOne(context.TODO(), r)
	if err != nil {
		return utils.Message(http.StatusNotModified, "Did not insert resume")
	}

	r.ID = insertResult.InsertedID.(primitive.ObjectID)

	response := utils.Message(http.StatusCreated, "Resume inserted")
	response["data"] = r

	return response
}

func GetResumes(userId primitive.ObjectID) map[string]interface{} {
	resumes, err := getResume(userId)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	response := utils.Message(http.StatusOK, "Resume found")
	response["data"] = resumes

	return response

}

func getResume(userId primitive.ObjectID) ([]*Resume, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	resume := db.Collection(ResumeCollection)
	result := make([]*Resume, 0)
	filter := bson.D{{"userid", userId}}

	cursor, err := resume.Find(context.TODO(), filter)

	if err != nil {
		return nil, fmt.Errorf("Error: %s", err.Error())
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var res Resume

		if err := cursor.Decode(&res); err != nil {
			return nil, fmt.Errorf("Error: %s", err.Error())
		}

		result = append(result, &res)
	}

	return result, nil

}

func GetResumeByID(id primitive.ObjectID) map[string]interface{} {
	res, err := getResumeById(id)

	if err != nil {
		return utils.Message(http.StatusNotFound, err.Error())
	}

	response := utils.Message(http.StatusOK, "Resume Found")
	response["data"] = res

	return response
}

func getResumeById(id primitive.ObjectID) (*Resume, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	resume := db.Collection("resume")
	var result Resume
	filter := bson.D{{"_id", id}}

	if err := resume.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

/************************************************* URL RECORD **********************************************************/

type URLRecordDAL interface {
	CreateURLRecord(u *URLRecord) map[string]interface{}
	GenerateRandomURL(userId primitive.ObjectID) map[string]interface{}
}

type URLRecord struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	URL       string             `json:"url"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	User      primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
}

func CreateURLRecord(u *URLRecord) map[string]interface{} {
	if u.URL == "" {
		return utils.Message(http.StatusBadRequest, "Illegal URL")
	}

	u.CreatedAt = time.Now().UTC()

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

func GenerateRandomURL(userId primitive.ObjectID) map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could Not Get Handle On DB")
	}

	records := db.Collection("records")
	record := &URLRecord{
		User: userId,
	}

	if err := generateAttempt(records, 0, 3, record); err != nil {
		return utils.Message(http.StatusConflict, "Could Not Generate a URL")
	}

	return CreateURLRecord(record)
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
