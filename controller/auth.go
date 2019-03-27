package controller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bkrmdahal/go_restAPI/model"
	"github.com/bkrmdahal/go_restAPI/services"
	"github.com/bkrmdahal/go_restAPI/utils"
	"github.com/mongodb/mongo-go-driver/bson"
)

// var
var cryto = services.Crypto{}

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
	tokenStruct := &model.Token{}
	token := services.CreateToken(user.UserId)
	tokenStruct.Token = token
	tokenStruct.UserId = uuids
	tokenStruct.CreatedAt = time.Now().Local()

	// save token
	insertResult, err := services.Db_token.InsertOne(context.TODO(), tokenStruct)
	utils.Log.Info("Inserted a user: ", insertResult.InsertedID)

	// hash the password
	hash, _ := cryto.Generate(user.Password)
	user.Password = hash

	insertResultUser, err := services.Db_user.InsertOne(context.TODO(), user)
	utils.Log.Info("Inserted a user: ", insertResultUser.InsertedID)

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

// signup the user
func Login(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		panic(err)
	}

	// search for user
	filters := bson.D{{"username", user.Username}}

	// get the result
	var userDB model.User
	errs := services.Db_user.FindOne(context.TODO(), filters).Decode(&userDB)
	if errs != nil {
		panic(err)
	}

	utils.Log.Info("Found a single document: %+v\n", userDB)

	userJson, err := json.Marshal(userDB)

	if err != nil {
		panic(err)
	}

	// compare the password
	errPassword := cryto.Compare(userDB.Password, user.Password)

	if errPassword != nil {
		panic("wrong password")
	}

	// make token
	tokenStruct := &model.Token{}
	token := services.CreateToken(userDB.UserId)
	tokenStruct.Token = token
	tokenStruct.UserId = userDB.UserId
	tokenStruct.CreatedAt = time.Now().Local()

	// save token
	insertResult, err := services.Db_token.InsertOne(context.TODO(), tokenStruct)
	utils.Log.Info("Inserted a user: ", insertResult.InsertedID)

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

// signup the user
func User(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("badrequest"))
	}

	// Get the JWT string from the cookie
	token := c.Value
	utils.Log.Info("cookies ", token)

	// search for user
	filters := bson.D{{"token", token}}

	// get the result
	var tokenDB model.Token
	errs := services.Db_token.FindOne(context.TODO(), filters).Decode(&tokenDB)
	utils.Log.Info("token detail  ", tokenDB.Token)

	if errs != nil {
		panic(err)
	}

	// tokenJson, err := json.Marshal(tokenDB)

	// if err != nil {
	// 	panic(err)
	// }
	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// search for user
	filterUser := bson.D{{"userid", tokenDB.UserId}}

	// get the result
	userDB := &model.User{}
	errUser := services.Db_user.FindOne(context.TODO(), filterUser).Decode(&userDB)
	if errUser != nil {
		log.Fatal(err)
	}

	utils.Log.Info("Found a single document: %+v\n", userDB)

	userDB.Password = ""
	userJson, err := json.Marshal(userDB)

	if err != nil {
		panic(err)
	}

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
