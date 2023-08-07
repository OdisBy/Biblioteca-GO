package local

type Book struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CoverUrl  string `json:"coverUrl"`
	Completed bool   `json:"completed"`
}

type BooksData struct {
	Books []Book `json:"items"`
}
