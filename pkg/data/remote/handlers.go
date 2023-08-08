package remote

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

const (
	baseUrlConst = "https://www.googleapis.com/books/v1/volumes"
)

func SearchBookHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Entrando no search book handler")
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	nameRequest := params["string"]

	println(nameRequest)

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
		fmt.Println("Erro no get api google")
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
