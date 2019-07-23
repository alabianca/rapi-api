package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func InitDB() {
	e := godotenv.Load()

	if e != nil {
		fmt.Printf("%s\n", e.Error())
	}

	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	user := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	mongoURI := os.Getenv("MONGODB_PROD") // will be set by heroku

	log.Printf("DB Host @ <%s>\n", host)
	log.Printf("DB Port @ <%s>\n", port)
	var dbURI string
	if mongoURI != "" {
		dbURI = mongoURI
	} else {
		dbURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, dbName)
	}

	clientOptions := options.Client().ApplyURI(dbURI)
	clientOptions.SetConnectTimeout(time.Duration(5 * time.Second))

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
}

// GetDB returns the db specified by the db_name env variable
func GetDB() (*mongo.Database, error) {
	if db == nil {
		return nil, fmt.Errorf("Not Connected to Database")
	}

	dbName := os.Getenv("db_name")

	return db.Database(dbName), nil
}
