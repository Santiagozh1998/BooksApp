package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "booksapp"
	password = "santiago1998"
)

var database *gorm.DB
var err error

//STRUCT
type Book struct {
	gorm.Model
	Title  string
	Author string
}

func ConnectDatabase() {

	databaseInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		host, port, user, dbname, password)

	database, err = gorm.Open("postgres", databaseInfo)

	if err != nil {
		println(err.Error())
	}

	database.AutoMigrate(&Book{})

	println("Connect with PostgreSQL")
}

func GetAllBooks(response http.ResponseWriter, request *http.Request) {
	var books []Book
	database.Find(&books)
	json.NewEncoder(response).Encode(&books)
}

func UpdateBook(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	var book Book
	var updatedBook Book
	database.Find(&book, params["id"])
	json.NewDecoder(request.Body).Decode(&updatedBook)
	database.Model(&book).Update(&updatedBook)
	json.NewEncoder(response).Encode(&book)
}

func CreateBook(response http.ResponseWriter, request *http.Request) {
	var book Book
	json.NewDecoder(request.Body).Decode(&book)
	database.Create(&book)
	json.NewEncoder(response).Encode(&book)
}

func DeleteBook(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	var book Book
	database.Find(&book, params["id"])
	database.Delete(&book)
	json.NewEncoder(response).Encode(&book)
}
