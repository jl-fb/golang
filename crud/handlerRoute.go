package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	staticDir = "/assets/"
)

//func getPeopleEndpoint(response http.ResponseWriter, request *http.Request) {}
//func getPersonEndpoint(response http.ResponseWriter, request *http.Request) {}

func handlerRoutes() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/person", applicationJSON(createPersonEndpoint)).Methods("POST")
	router.HandleFunc("/people", applicationJSON(getPeopleEndpoint)).Methods("GET")
	router.HandleFunc("/person/{id}", applicationJSON(getPersonEndpoint)).Methods("GET")
	router.HandleFunc("/person/{id}", applicationJSON(updatePersonEndPoint)).Methods("PUT")
	router.HandleFunc("/person/{id}", applicationJSON(deletePersonEndPoint)).Methods("DELETE")
	//sempre por ultimo possivelmente por fazer a aplicação parar na "/"
	// ou usar "/home" que fica ok
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("."+staticDir))))
	http.ListenAndServe(":12345", router)
}
