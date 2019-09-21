package models

// Book is datamodel for a book info
type Book struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Lang    string    `json:"lang"`
	Intro   string    `json:"intro"`
	Cover   string    `json:"cover"`
	Toc     []Chapter `json:"toc"`
	Authors []Person  `json:"authors"`
}
