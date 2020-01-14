package main

import (
	"log"
	"net/http"

	"tiny_server/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/user", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", handlers.RetrieveUsers).Methods("GET")

	log.Println("Server was started...")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
