package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func init() {
	e := godotenv.Load()

	if e != nil {
		fmt.Print(e)
	}

	host := os.Getenv("db_host")
	port := os.Getenv("db_port")

	dbURI := fmt.Sprintf("mongodb://%s:%s", host, port)

	clientOptions := options.Client().ApplyURI(dbURI)

	// connect to mongodb
	var err error
	db, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// check connection
	err = db.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to MongoDB @ %s\n", dbURI)
}

func GetDB() (*mongo.Database, error) {
	if db == nil {
		return nil, fmt.Errorf("Not Connected to Database")
	}

	dbName := os.Getenv("db_name")

	return db.Database(dbName), nil
}
