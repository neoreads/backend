package models

type Content struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Chapter struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}
