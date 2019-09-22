package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// define our database connection
var db *gorm.DB

// Visitors Struct defines the structure to be used for the Visitor table
type Visitors struct {
	gorm.Model
	Name  string
	Email string
}

// function to initialize DB connection before we even create the tables
func initialMigration() {
	db, err := gorm.Open("sqlite3", "gormTestDatabase.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// actual migration i.e. actual creation of the tables
	db.AutoMigrate(&Visitors{})
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Visitors{})
}

// abstract away the database connection code, to make things DRY
func databaseConnection() (db *gorm.DB) {
	db, err := gorm.Open("sqlite3", "gormTestDatabase.db")
	if err != nil {
		panic(err.Error())
	}
	// notice we're not calling a defer db.Close() here as this is not where defer is needed
	return db
}

func getAllVisitorsHandler(w http.ResponseWriter, r *http.Request) {
	// db, err := gorm.Open("sqlite3", "gormTestDatabase.db")
	// if err != nil {
	//	panic(err.Error())
	//}
	db := databaseConnection()
	defer db.Close()

	var visitors []Visitors             // create a slice of type Visitors i.e. an array of elements of type struct
	db.Find(&visitors)                  // returns all records in table "Visitor" and saves to the slice "Visitors"
	json.NewEncoder(w).Encode(visitors) // convert to JSON and pass the encoded result to the responsewriter `w`
}

func createVisitorHandler(w http.ResponseWriter, r *http.Request) {
	db := databaseConnection()
	defer db.Close()

	vars := mux.Vars(r)
	requestParameterName := vars["name"]
	requestParameterEmail := vars["email"]

	db.Create(&Visitors{Name: requestParameterName, Email: requestParameterEmail}) // insert data into database using the Struct itself
	dataToRender := "New Visitor successfully created"
	io.WriteString(w, dataToRender)
}

func deleteVisitorHandler(w http.ResponseWriter, r *http.Request) {
	db := databaseConnection()
	defer db.Close()

	vars := mux.Vars(r)
	requestParameterName := vars["name"]

	var visitor Visitors                                      // initialize variable to response data of what we want to delete
	db.Where("name = ?", requestParameterName).Find(&visitor) // find what to delete
	db.Delete(&visitor)

	dataToRender := "Successfully deleted Visitor"
	io.WriteString(w, dataToRender)
}

func updateVisitorHandler(w http.ResponseWriter, r *http.Request) {
	db := databaseConnection()
	defer db.Close()

	vars := mux.Vars(r)
	requestParameterName := vars["name"]
	requestParameterEmail := vars["email"]

	var visitor Visitors
	db.Where("name = ?", requestParameterName).Find(&visitor) // find visitor to update
	visitor.Email = requestParameterEmail                     // assign the new email to the visitor's email attribute
	db.Save(&visitor)                                         // save the changes to the visitor object

	dataToRender := "Successfully updated Visitor details"
	io.WriteString(w, dataToRender)
}
