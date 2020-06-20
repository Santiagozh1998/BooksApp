package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Santiagozh1998/BooksApp/database"
	"github.com/gorilla/mux"
)

func main() {

	//PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	//DATABASE
	database.ConnectDatabase()

	//ROUTES
	route := mux.NewRouter()
	staticFiles := http.FileServer(http.Dir("views/assets/"))
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFiles))

	route.Handle("/", http.FileServer(http.Dir("./views")))
	route.HandleFunc("/api/books/all", database.GetAllBooks).Methods("GET")
	route.HandleFunc("/api/books/create", database.CreateBook).Methods("POST")
	route.HandleFunc("/api/books/update/{id}", database.UpdateBook).Methods("PUT")
	route.HandleFunc("/api/books/delete/{id}", database.DeleteBook).Methods("DELETE")

	//STARTING SERVER
	log.Fatal(http.ListenAndServe(":"+port, route))
}
