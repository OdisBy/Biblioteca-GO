package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type Book struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CoverUrl  string `json:"coverUrl"`
	Completed bool   `json:"completed"`
}

func getBooks(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(writer, "Method is not supported", http.StatusBadRequest)
		return
	}
	tmpl, err := template.ParseGlob(filepath.Join("./static", "*.html"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(writer, "index.html", books)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, item := range books {
		if item.ID == params["id"] {
			err := json.NewEncoder(writer).Encode(item)
			if err != nil {
				return
			}
			return
		}
	}
}

func createBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var newBook Book
	err := json.NewDecoder(request.Body).Decode(&newBook)
	if err != nil {
		return
	}
	newBook.ID = strconv.Itoa(len(books) + 1)
	books = append(books, newBook)
	err = json.NewEncoder(writer).Encode(newBook)
	if err != nil {
		return
	}
}

func updateBook(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

}

func deleteBook(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

}

func completeBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	for i, item := range books {
		if item.ID == params["id"] {
			var newBook Book
			err := json.NewDecoder(request.Body).Decode(&newBook)
			if err != nil {
				return
			}
			books[i].Completed = newBook.Completed
			books = append(books, newBook)
			err = json.NewEncoder(writer).Encode(newBook)
			if err != nil {
				return
			}
			log.Printf("\nComplete book %s, book name: %s, book cover: %s", books[i].ID, books[i].Name, books[i].CoverUrl)
			return
		}
	}
}

func defaultHandler(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Missing path", http.StatusBadRequest)
	return
}

var books []Book

func main() {
	books = append(books,
		Book{ID: "1", CoverUrl: "https://i.imgur.com/EZdSNjv.jpeg", Name: "Harry Potter: A Câmera Secreta", Completed: false},
		Book{ID: "2", CoverUrl: "https://i.imgur.com/BwtS9oD.png", Name: "Teste 2", Completed: true},
	)

	//fileServer := http.FileServer(http.Dir("./static"))

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("POST")
	router.HandleFunc("/books/complete/{id}", completeBook).Methods("POST")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	//router.PathPrefix("/").Handler(fileServer)
	router.HandleFunc("/", getBooks)

	fmt.Printf("Starting listening at localhost:5000")

	log.Fatal(http.ListenAndServe(":5000", router))
}
