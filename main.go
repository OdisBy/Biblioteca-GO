package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Book struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CoverUrl  string `json:"coverUrl"`
	Completed bool   `json:"completed"`
}

type BooksData struct {
	Books []Book `json:"books"`
}

var booksData BooksData

func getNextID() string {
	if len(booksData.Books) == 0 {
		return "1"
	}

	lastID, err := strconv.Atoi(booksData.Books[len(booksData.Books)-1].ID)
	if err != nil {
		return "1"
	}

	nextID := lastID + 1
	return strconv.Itoa(nextID)
}

func loadBooksData() error {
	file, err := os.OpenFile("./data/books.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		booksData = BooksData{
			Books: []Book{},
		}
		return nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&booksData)
	if err != nil {
		return err
	}

	return nil
}

func saveBooksData() error {
	file, err := os.Create("./data/books.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(booksData)
	if err != nil {
		return err
	}

	return nil
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
	err = tmpl.ExecuteTemplate(writer, "index.html", booksData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, item := range booksData.Books {
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

	err := loadBooksData()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var newBook Book
	err = json.NewDecoder(request.Body).Decode(&newBook)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	newBook.ID = getNextID()

	booksData.Books = append(booksData.Books, newBook)

	err = saveBooksData()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(newBook)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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

	for i, item := range booksData.Books {
		if item.ID == params["id"] {
			var newBook Book
			err := json.NewDecoder(request.Body).Decode(&newBook)
			if err != nil {
				return
			}
			booksData.Books[i].Completed = newBook.Completed
			err = saveBooksData()
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(writer).Encode(booksData.Books[i])
			if err != nil {
				return
			}
			log.Printf("\nComplete book %s, book name: %s, book cover: %s", booksData.Books[i].ID, booksData.Books[i].Name, booksData.Books[i].CoverUrl)
			return
		}
	}
}

func defaultHandler(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Missing path", http.StatusBadRequest)
	return
}

func main() {
	err := loadBooksData()
	if err != nil {
		log.Fatal("Erro ao carregar os dados:", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("POST")
	router.HandleFunc("/books/complete/{id}", completeBook).Methods("POST")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/", getBooks).Methods("GET")

	fmt.Printf("Starting listening at localhost:5000")

	log.Fatal(http.ListenAndServe(":5000", router))
}
