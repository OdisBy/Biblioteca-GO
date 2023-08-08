package local

import (
	"encoding/json"
	"log"
	"os"
)

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	CoverUrl    string `json:"coverUrl"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type BooksData struct {
	Books []Book `json:"books"`
}

func LoadBooksData(target *BooksData) error {
	file, err := os.OpenFile("./pkg/data/local/books.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		*target = BooksData{
			Books: []Book{},
		}
		return nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func SaveBooksData(target *BooksData) error {
	file, err := os.Create("./data/books.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(target)
	if err != nil {
		return err
	}
	return nil
}
