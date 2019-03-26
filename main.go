package main

import (
	"net/http"

	"github.com/bkrmdahal/go_restAPI/controller"
	"github.com/bkrmdahal/go_restAPI/utils"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello worlds!"))

}

func main() {
	// start the logger
	utils.Log.Info("Starting the server")

	mux := http.NewServeMux()

	mux.HandleFunc("/signup", controller.Signup)
	mux.HandleFunc("/welcome", welcome)

	// start the server at 5000
	http.ListenAndServe(":5000", mux)
	utils.Log.Info("stopping  the server")

}
