package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/alabianca/rapi-api/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

type TokenDAL interface {
	GetToken(email, password string) map[string]interface{}
}

type Token struct {
	UserID      string `json:"userId"`
	TokenString string `json:"token"`
	Expires     string `json:"expires"`
	jwt.StandardClaims
}

type TokenSource struct{}

func (ts TokenSource) GetToken(email, password string) map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on database")
	}

	user := &User{}
	filter := bson.D{{"email", email}}
	users := db.Collection(UserCollection)
	expiresInSeconds, _ := strconv.ParseInt(os.Getenv("token_expiry"), 10, 64)
	log.Printf("Token expires in %d seconds\n", expiresInSeconds)
	expires := time.Now().Add(time.Duration(int64(time.Second) * expiresInSeconds))
	if err := users.FindOne(context.TODO(), filter).Decode(user); err == mongo.ErrNoDocuments {
		return utils.Message(http.StatusNotFound, fmt.Sprintf("User %s not found", email))
	}

	// verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return utils.Message(http.StatusUnauthorized, "Authentication Error")
	}

	// user is legit. send up a token
	tk := &Token{UserID: user.ID.Hex()}
	//tk.ExpiresAt = time.Now().UnixNano() + (expiresInSeconds * int64(time.Nanosecond))
	tk.ExpiresAt = expires.Unix()

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))

	if err != nil {
		return utils.Message(http.StatusUnauthorized, "Token Error "+err.Error())
	}

	tk.TokenString = tokenString
	tk.Expires = expires.UTC().String()

	response := utils.Message(http.StatusOK, "Successfully got a token")
	response["data"] = tk

	return response
}
