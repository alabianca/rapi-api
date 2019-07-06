package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alabianca/rapi-api/controllers"
)

func main() {
	api := controllers.API{
		DAL: controllers.DefaultDAL{},
	}
	router := apiRoutes(&api)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	if host == "localhost" {
		host = ""
	}

	address := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server Listening @ %s\n", address)

	err := http.ListenAndServe(address, router)

	if err != nil {
		log.Println(err)
	}
}
