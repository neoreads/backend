package repositories

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	pmodels "github.com/neoreads/backend/prepare/models"
	"github.com/neoreads/backend/server/models"
)

// BookRepo book related data repository
type BookRepo struct {
	db      *sqlx.DB
	rootDir string
}

// NewBookRepo creator for BookRepo
func NewBookRepo(db *sqlx.DB, root string) *BookRepo {
	return &BookRepo{db: db, rootDir: root}
}

// GetBook get book info by bookid
func (r *BookRepo) GetBook(id string) (book models.Book, found bool) {
	err := r.db.Get(&book, "SELECT * from books where id=$1", id)
	if err != nil {
		return models.Book{}, false
	}
	return book, true
}

// ContainsBookID check if bookid is already contained in the table
func (r *BookRepo) ContainsBookID(bookid string) bool {
	count := 0
	err := r.db.Get(&count, "SELECT count(id) from books where id=$1", bookid)
	if err != nil {
		log.Fatalf("[BookRepo] error checking book id: %s", err)
		return true
	}
	return count > 0
}

// GetTOC get table of contents
func (r *BookRepo) GetTOC(id string) (toc []models.Chapter) {
	err := r.db.Select(&toc, "SELECT * from chapters where bookid=$1 order by \"order\" asc", id)
	if err != nil {
		log.Printf("err:%s\n", err)
	}
	return toc
}

func (r *BookRepo) readText(chap *models.Chapter) string {
	bookid := chap.BookID
	chapid := chap.ID
	dir := bookid[:4]
	path := r.rootDir + "books/" + dir + "/" + bookid + "/" + chapid + ".md"
	log.Printf("reading chapter from %s\n", path)
	text, err := ioutil.ReadFile(path)
	if err == nil {
		return string(text)
	}
	return ""
}

// GetContent get the content of a chapter by bookid and chapid
// Note: chapid in database is actually bookid+chapid
func (r *BookRepo) GetContent(bookid string, chapid string) (content models.Content, found bool) {
	log.Printf("Getting content:%s:%s", bookid, chapid)
	chap := &models.Chapter{}
	log.Printf("SELECT * from chapters where bookid=%s and id=%s\n", bookid, chapid)
	err := r.db.Get(chap, "SELECT * from chapters where bookid=$1 and id=$2", bookid, chapid)
	log.Println(chap)
	if err == nil {
		text := r.readText(chap)
		content.Content = text
		content.ID = chap.ID
		content.Title = chap.Title
		return content, true
	}
	log.Printf("erro:%s", err)
	return content, false
}

func (r *BookRepo) AddBook(toc *pmodels.Toc) {
	r.db.Exec("INSERT INTO books VALUES($1, $2)", toc.BookID, toc.Title)
	for od, item := range toc.Items {
		_, err := r.db.Exec("INSERT INTO chapters VALUES($1, $2, $3, $4)",
			item.ChapID, od+1, toc.BookID, item.Title)
		if err != nil {
			log.Fatal(err)
		}
	}
}
