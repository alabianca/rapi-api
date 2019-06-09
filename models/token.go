package models

import (
	"context"
	"fmt"
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

type Token struct {
	UserID      string `json:"userId"`
	TokenString string `json:"token"`
	Expires     string `json:"expires"`
	jwt.StandardClaims
}

func GetToken(email, password string) map[string]interface{} {
	db, err := GetDB()

	if err != nil {
		return utils.Message(http.StatusInternalServerError, "Could not get a handle on database")
	}

	user := &User{}
	filter := bson.D{{"email", email}}
	users := db.Collection("users")
	expiresInSeconds, _ := strconv.ParseInt(os.Getenv("token_expiry"), 10, 16)
	expires := time.Now().Add(time.Duration(int64(time.Nanosecond) * expiresInSeconds))
	if err := users.FindOne(context.TODO(), filter).Decode(user); err == mongo.ErrNoDocuments {
		return utils.Message(http.StatusNotFound, fmt.Sprintf("User %s not found", email))
	}

	// verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return utils.Message(http.StatusForbidden, "Authentication Error")
	}

	// user is legit. send up a token
	tk := &Token{UserID: user.ID.Hex()}
	// tk.ExpiresAt = time.Now().Add(time.Second * expires)
	tk.ExpiresAt = time.Now().UnixNano() + (expiresInSeconds * int64(time.Nanosecond))

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))

	if err != nil {
		return utils.Message(http.StatusForbidden, "Token Error "+err.Error())
	}

	tk.TokenString = tokenString
	tk.Expires = expires.UTC().String()

	response := utils.Message(http.StatusOK, "Successfully got a token")
	response["token"] = tk

	return response
}
