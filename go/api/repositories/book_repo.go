package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/neoreads-backend/go/api/datamodels"
)

type BookRepo interface {
	Get(id string) (book datamodels.Book, found bool)
}

type bookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) BookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) Get(id string) (book datamodels.Book, found bool) {
	err := r.db.Get(&book, "SELECT * from book where id=$1", id)
	if err != nil {
		return datamodels.Book{}, false
	}
	return book, true
}
