package main

import (
	"fmt"
	"github.com/OdisBy/Biblioteca-GO/pkg/data/local"
	"github.com/OdisBy/Biblioteca-GO/pkg/data/remote"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/books", local.GetBooks).Methods("GET")
	router.HandleFunc("/book/{id}", local.GetBooks).Methods("GET")
	router.HandleFunc("/book", local.CreateBook).Methods("POST")
	router.HandleFunc("/book/{id}", local.UpdateBook).Methods("POST")
	router.HandleFunc("/book/complete/{id}", local.CompleteBook).Methods("POST")
	router.HandleFunc("/book/{id}", local.DeleteBook).Methods("DELETE")
	router.HandleFunc("/search/{string}", remote.SearchBookHandler).Methods("GET")
	router.HandleFunc("/", local.GetBooks).Methods("GET")

	fmt.Printf("Starting listening at localhost:5000")

	log.Fatal(http.ListenAndServe(":5000", router))
}
