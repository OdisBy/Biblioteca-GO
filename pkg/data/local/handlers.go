package local

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var localBooksData BooksData

func getNextID() string {
	if len(localBooksData.Books) == 0 {
		return "1"
	}

	lastID, err := strconv.Atoi(localBooksData.Books[len(localBooksData.Books)-1].ID)
	if err != nil {
		return "1"
	}

	nextID := lastID + 1
	return strconv.Itoa(nextID)
}

func GetBooks(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(writer, "Method is not supported", http.StatusBadRequest)
		return
	}
	err := LoadBooksData(&localBooksData)
	if err != nil {
		return
	}
	tmpl, err := template.ParseGlob(filepath.Join("./static", "*.html"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(writer, "index.html", localBooksData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetBook(writer http.ResponseWriter, request *http.Request) {
	err := LoadBooksData(&localBooksData)
	if err != nil {
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, item := range localBooksData.Books {
		if item.ID == params["id"] {
			err := json.NewEncoder(writer).Encode(item)
			if err != nil {
				return
			}
			return
		}
	}
}

func CreateBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	err := LoadBooksData(&localBooksData)
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

	localBooksData.Books = append(localBooksData.Books, newBook)

	err = SaveBooksData(&localBooksData)
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

func UpdateBook(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

}

func DeleteBook(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

}

func CompleteBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	err := LoadBooksData(&localBooksData)
	if err != nil {
		return
	}

	for i, item := range localBooksData.Books {
		if item.ID == params["id"] {
			var newBook Book
			err := json.NewDecoder(request.Body).Decode(&newBook)
			if err != nil {
				return
			}
			localBooksData.Books[i].Completed = newBook.Completed
			err = SaveBooksData(&localBooksData)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(writer).Encode(localBooksData.Books[i])
			if err != nil {
				return
			}
			log.Printf("\nComplete book %s, book name: %s, book cover: %s", localBooksData.Books[i].ID, localBooksData.Books[i].Title, localBooksData.Books[i].CoverUrl)
			return
		}
	}
}

func SearchBookHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	nameRequest := params["string"]

	var results []Book

	if len(nameRequest) >= 3 {
		for _, book := range localBooksData.Books {
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(nameRequest)) {
				results = append(results, book)
			}
		}
	}

	err := json.NewEncoder(writer).Encode(results)
	if err != nil {
		return
	}
}
