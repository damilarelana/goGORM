package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	dataToRender := "Just the beginning of GORM thingies"
	io.WriteString(w, dataToRender)
}

func serviceRequestHandlers() {
	newRouter := mux.NewRouter().StrictSlash(true)
	newRouter.HandleFunc("/visitors", getAllVisitorsHandler).Methods("GET")
	newRouter.HandleFunc("/visitor/{name}/{email}", createVisitorHandler).Methods("POST")
	newRouter.HandleFunc("/visitor/{name}", deleteVisitorHandler).Methods("DELETE")
	newRouter.HandleFunc("/visitor/{name}/{email}", updateVisitorHandler).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8081", newRouter))
}

func main() {
	fmt.Println("a little bit of GORM")

	// create the database tables needed by the APIs
	go initialMigration()

	// start the API service
	go serviceRequestHandlers()

	// create an artificial pause "to ensure the main function goroutine does not cause the serviceRequestHandler goroutine to exit"
	var tempString string
	fmt.Scanln(&tempString)
}
