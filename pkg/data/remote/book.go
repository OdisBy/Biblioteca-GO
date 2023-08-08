package remote

type SearchBooksData struct {
	Books []SearchBook `json:"items"`
}

type SearchBook struct {
	ID         string           `json:"id"`
	VolumeInfo SearchVolumeInfo `json:"volumeInfo"`
}

type Book struct {
	ID         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type SearchVolumeInfo struct {
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Description string   `json:"description"`
}

type VolumeInfo struct {
	Title         string     `json:"title"`
	Authors       []string   `json:"authors"`
	PublishedDate string     `json:"publishedDate"`
	Description   string     `json:"description"`
	PageCount     int        `json:"pageCount"`
	Categories    []string   `json:"categories"`
	ImageLinks    ImageLinks `json:"imageLinks"`
	Completed     bool       `json:"completed"`
}

type ImageLinks struct {
	Small string `json:"thumbnail"`
}
