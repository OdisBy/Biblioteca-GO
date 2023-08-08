package remote

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
)

const (
	baseUrlConst = "https://www.googleapis.com/books/v1/volumes"
)

func SearchBookHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	nameRequest := params["string"]

	var searchBooks SearchBooksData

	baseUrl, _ := url.Parse(baseUrlConst)
	urlParams := url.Values{}
	urlParams.Add("q", nameRequest)
	urlParams.Add("maxResults", "3")
	urlParams.Add("langRestrict", "pt")
	urlParams.Add("printType", "books")
	baseUrl.RawQuery = urlParams.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		fmt.Println("Erro no get by name api google")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&searchBooks)
	if err != nil {
		fmt.Printf("Erro ao decodificar resposta JSON: %v\n", err)
		return
	}

	err = json.NewEncoder(writer).Encode(searchBooks)
	if err != nil {
		return
	}
}

func GetBookInfo(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	idRequest := params["id"]

	var book Book
	baseUrl, _ := url.Parse(baseUrlConst + "/" + idRequest)
	urlParams := url.Values{}
	urlParams.Add("langRestrict", "pt")
	urlParams.Add("printType", "books")
	baseUrl.RawQuery = urlParams.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		fmt.Println("Erro no get book by id api google")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&book)
	if err != nil {
		fmt.Printf("Erro ao decodificar resposta JSON: %v\n", err)
		return
	}

	if book.VolumeInfo.ImageLinks.Small == "" {
		book.VolumeInfo.ImageLinks.Small = "../static/default_cover.jpg"
	}

	tmpl, err := template.ParseGlob(filepath.Join("./static", "*.html"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(writer, "bookPage.html", book)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
