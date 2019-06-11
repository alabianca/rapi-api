.PHONY:
server: deps
	go build -o bin/rapid

deps:
	go get github.com/joho/godotenv
	go get github.com/dgrijalva/jwt-go
	go get github.com/go-chi/chi
	go get github.com/go-chi/cors
	go get github.com/go-chi/chi/middleware
	go get github.com/go-chi/render
	go get go.mongodb.org/mongo-driver