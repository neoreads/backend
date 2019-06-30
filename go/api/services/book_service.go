package services

import (
	"github.com/neoreads-backend/go/api/datamodels"
	"github.com/neoreads-backend/go/api/repositories"
)

// BookService handles CRUD and other operations for datamodel Book
type BookService interface {
	GetByID(id string) (datamodels.Book, bool)
}

type bookService struct {
	repo repositories.BookRepo
}

// NewBookService returns a new bookService instance
func NewBookService(repo repositories.BookRepo) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) GetByID(id string) (datamodels.Book, bool) {
	return s.repo.Get(id)
}
