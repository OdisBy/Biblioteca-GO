package remote

type BookApi struct {
	ID         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type VolumeInfo struct {
	Title         string     `json:"title"`
	Authors       Authors    `json:"authors"`
	PublishedDate string     `json:"publishedDate"`
	Description   string     `json:"description"`
	Categories    Categories `json:"categories"`
	CoverUrl      string     `json:"coverUrl"`
	Completed     bool       `json:"completed"`
}

type Authors struct {
	names []string
}

type Categories struct {
	categories []string
}

type BooksDataApi struct {
	Books []BookApi `json:"books"`
}
