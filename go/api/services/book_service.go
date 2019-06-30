package services

import (
	"github.com/neoreads-backend/go/api/datamodels"
)

// BookService handles CRUD and other operations for datamodel Book
type BookService interface {
	GetByID(id string) (datamodels.Book, bool)
}

type bookService struct {
}

// NewBookService returns a new bookService instance
func NewBookService() BookService {
	return &bookService{}
}

func (s *bookService) GetByID(id string) (datamodels.Book, bool) {
	return datamodels.Book{
		ID:     id,
		Title:  "To Kill a Mocking Bird",
		Author: "Harper Lee",
		Desc:   "....",
	}, true
}
