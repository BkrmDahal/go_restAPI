package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bkrmdahal/go_restAPI/model"
	"github.com/bkrmdahal/go_restAPI/services"
	"go.mongodb.org/mongo-driver/bson"
)

// signup the user
func Login_required(w http.ResponseWriter, r *http.Request) {
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

	// search for user
	filters := bson.D{{"token", token}}

	// get the result
	var tokenDB model.Token
	errs := services.Db_token.FindOne(context.TODO(), filters).Decode(&tokenDB)

	if errs != nil {
		panic(err)
	}

	tokenJson, err := json.Marshal(tokenDB)

	if err != nil {
		panic(err)
	}
	fmt.Println("hello test", tokenJson)

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	// 	if time.Unix(token.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	} else {
	// 		fmt.Println("hello test")
	// 	}
}
