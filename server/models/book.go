package models

// Book is datamodel for a book info
type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Intro string `json:"intro"`
	Cover string `json:"cover"`
}
