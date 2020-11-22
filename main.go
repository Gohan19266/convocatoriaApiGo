package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	s := router.PathPrefix("/convocatorias").Subrouter()
	s.HandleFunc("/all", getAllConvocatorias).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", s))
}
