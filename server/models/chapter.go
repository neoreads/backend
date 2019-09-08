package models

type Content struct {
	ID      string `json:"id"`
	BookID string `json:"bookid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Chapter struct {
	ID     string `json:"id"`
	Order  int    `json:"order"`
	BookID string `json:"bookid"`
	Title  string `json:"title"`
	Content string `json:"content"`
}
