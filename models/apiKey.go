package models

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/alabianca/rapi-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const KeyCollection = "keys"

type APIKey struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"userId,omitempty"`
	Resume       primitive.ObjectID `json:"resume,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	Key          string             `json:"key"`
	Scope        []string           `json:"scope"`
	FriendlyName string             `json:"friendlyName"`
}

func (a *APIKey) Create() map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get handle on db")
	}

	keys := db.Collection(KeyCollection)
	a.Key = makeKey()
	a.CreatedAt = time.Now().UTC()

	res, err := keys.InsertOne(context.TODO(), a)

	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	a.ID = res.InsertedID.(primitive.ObjectID)

	response := utils.Message(http.StatusCreated, "API Key Created")
	response["data"] = a

	return response

}

func (a *APIKey) UpdateKey() map[string]interface{} {
	db, err := GetDB()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on db")
	}

	keys := db.Collection(KeyCollection)

	filter := bson.D{{"_id", a.ID}}
	update := bson.D{
		{"$set", bson.D{{"scope", a.Scope}}},
		{"$set", bson.D{{"friendlyname", a.FriendlyName}}},
	}

	_, err = keys.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Error Occurerred %s", err.Error())
		return utils.Message(http.StatusNotModified, err.Error())
	}

	response := utils.Message(http.StatusOK, "Key Updated successfully")
	response["data"] = a

	return response

}

func DeleteKey(id primitive.ObjectID) map[string]interface{} {
	db, err := GetDB()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on db")
	}

	keys := db.Collection(KeyCollection)

	filter := bson.D{{"_id", id}}
	_, err = keys.DeleteOne(context.TODO(), filter)
	if err != nil {
		return utils.Message(http.StatusNotModified, "Could not delete key "+err.Error())
	}

	response := utils.Message(http.StatusOK, "Key Deleted")
	response["data"] = "OK"
	return response
}

func GetKeys(userID, resumeID primitive.ObjectID) map[string]interface{} {
	db, err := GetDB()
	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on db")
	}

	keys := db.Collection(KeyCollection)
	filter := bson.D{{"userid", userID}, {"resume", resumeID}}

	cursor, err := keys.Find(context.TODO(), filter)
	if err != nil {
		return utils.Message(http.StatusInternalServerError, err.Error())
	}

	results := make([]*APIKey, 0)

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var key APIKey
		if err := cursor.Decode(&key); err != nil {
			return utils.Message(http.StatusInternalServerError, err.Error())
		}

		results = append(results, &key)
	}

	response := utils.Message(http.StatusOK, "Found Keys")
	response["data"] = results

	return response

}

func makeKey() string {
	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	encoder.Write(randBytes(32))

	return buf.String()
}

func randBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}
