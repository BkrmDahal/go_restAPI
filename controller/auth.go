package controller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bkrmdahal/go_restAPI/model"
	"github.com/bkrmdahal/go_restAPI/services"
)

// signup the user
func Signup(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		panic(err)
	}

	user.CreatedAt = time.Now().Local()

	// make unique user id
	uuidb := make([]byte, 4) //equals 8 charachters
	rand.Read(uuidb)
	uuids := hex.EncodeToString(uuidb)
	user.UserId = uuids

	// make token
	token := services.CreateToken(user.Password)
	user.Token = token

	// hash the password
	cryto := services.Crypto{}
	hash, _ := cryto.Generate(user.Password)
	user.Password = hash

	insertResult, err := services.Db_user.InsertOne(context.TODO(), user)
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	userJson, err := json.Marshal(user)

	if err != nil {
		panic(err)
	}

	// set the token
	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	// make responses to send the request
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)

}
