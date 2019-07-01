package repositories

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads-backend/go/server/models"
)

// BookRepo book related data repository
type BookRepo struct {
	db *sqlx.DB
}

// NewBookRepo creator for BookRepo
func NewBookRepo(db *sqlx.DB) *BookRepo {
	return &BookRepo{db: db}
}

// GetBook get book info by bookid
func (r *BookRepo) GetBook(id string) (book models.Book, found bool) {
	err := r.db.Get(&book, "SELECT * from book where id=$1", id)
	if err != nil {
		return models.Book{}, false
	}
	return book, true
}

func readText(path string) string {
	file := "D:/neoreads/data/000/" + path
	log.Println(file)
	text, err := ioutil.ReadFile(file)
	if err == nil {
		return string(text)
	}
	return ""
}

// GetContent get the content of a chapter by bookid and chapid
// Note: chapid in database is actually bookid+chapid
func (r *BookRepo) GetContent(bookid string, chapid string) (content models.Content, found bool) {
	log.Println("Getting content")
	id := bookid + chapid
	chap := models.Chapter{}
	err := r.db.Get(&chap, "SELECT * from chapter where id=$1", id)
	log.Println(chap)
	if err == nil {
		path := chap.Path
		log.Println(path)
		text := readText(path)
		content.Content = text
		content.ID = chap.ID
		content.Title = chap.Title
		log.Print(err)
		return content, true
	}
	return content, false
}
