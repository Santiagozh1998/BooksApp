package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//STRUCT
type Book struct {
	ID     primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	Title  string             `json:"title, omitempty" bson:"title, omitempty"`
	Author string             `json:"author, omitempty" bson:"author, omitempty"`
}

//FUNCTIONS ROUTES
var client *mongo.Client

func getBooks(response http.ResponseWriter, request *http.Request) {

	println("getBooks")

	var books []*Book
	collection := client.Database("books").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for results.Next(ctx) {
		var book Book
		results.Decode(&book)
		books = append(books, &book)
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
	results.Close(ctx)
	json.NewEncoder(response).Encode(books)
}

func createBooks(response http.ResponseWriter, request *http.Request) {

	println("CreateBooks")

	var book Book
	json.NewDecoder(request.Body).Decode(&book)
	collection := client.Database("books").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results, err := collection.InsertOne(ctx, bson.M{"Title": book.Title, "Author": book.Author})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(results.InsertedID)

}

func updateBooks(response http.ResponseWriter, request *http.Request) {

	println("UpdateBooks")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	json.NewDecoder(request.Body).Decode(&book)
	collection := client.Database("books").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results, err := collection.ReplaceOne(ctx,
		bson.M{"_id": id},
		bson.M{"Title": book.Title, "Author": book.Author})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(results.ModifiedCount)
}

func deleteBooks(response http.ResponseWriter, request *http.Request) {

	println("DeleteBooks")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("books").Collection("book")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(results.DeletedCount)
}

//CREATING ROUTES
func newRouter() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/api/books/", getBooks).Methods("GET")
	r.HandleFunc("/api/books/", createBooks).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBooks).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBooks).Methods("DELETE")

	return r
}

func main() {

	var user = "Ingresar el usuario de MongoDB Atlas"
	var password = "Ingresar contraseña del usuario"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://"+user+":"+password+"@books-shard-00-00-favnn.mongodb.net:27017,books-shard-00-01-favnn.mongodb.net:27017,books-shard-00-02-favnn.mongodb.net:27017/books?ssl=true&replicaSet=Books-shard-0&authSource=admin&retryWrites=true&w=majority",
	))

	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	println("Connected to MongoDB")

	//ROUTES
	router := newRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	println("Starting server on port:" + port)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})

	//STARTING SERVER
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origin)(router)))
}
