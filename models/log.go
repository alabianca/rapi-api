package models

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alabianca/rapi-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const APILogsCollection = "apiLogs"

type LogDAL interface {
	CeateLog(l *Log) map[string]interface{}
	GetLogsFor(apiKeyID primitive.ObjectID) map[string]interface{}
	GetLogsForKeys(apiIds []primitive.ObjectID) map[string]interface{}
}

type Log struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	APIID primitive.ObjectID `json:"apiId"`
	Date  time.Time          `json:"date"`
}

func CreateLog(l *Log) map[string]interface{} {
	logs, err := getLogsCollection()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	res, err := logs.InsertOne(context.TODO(), l)
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	l.ID = res.InsertedID.(primitive.ObjectID)

	response := utils.Message(http.StatusCreated, "Created Log")
	response["data"] = l

	return response
}

func GetLogsFor(apiKeyID primitive.ObjectID) map[string]interface{} {
	logs, err := getLogsCollection()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	filter := bson.D{{"apiid", apiKeyID}}

	result, err := findLogs(logs, filter)
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	response := utils.Message(http.StatusOK, "Logs Found")
	response["data"] = result

	return response
}

func GetLogsForKeys(apiIds []primitive.ObjectID) map[string]interface{} {
	logs, err := getLogsCollection()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}
	log.Printf("Looking for %s\n", apiIds)

	filter := bson.D{
		{"apiid", bson.D{{
			"$in", apiIds,
		}}},
	}

	result, err := findLogs(logs, filter)
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	response := utils.Message(http.StatusOK, "Logs Found")
	response["data"] = result

	return response

}

func findLogs(logs *mongo.Collection, filter bson.D) ([]*Log, error) {
	result := make([]*Log, 0)

	cursor, err := logs.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {

		var logEvent Log
		if err := cursor.Decode(&logEvent); err != nil {
			return nil, err
		}

		result = append(result, &logEvent)
	}

	return result, nil

}

func getLogsCollection() (*mongo.Collection, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	return db.Collection(APILogsCollection), nil
}
