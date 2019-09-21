package repositories

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path"

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
	toc := r.GetTOC(id)
	book.Toc = toc
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
		content.BookID = bookid
		return content, true
	}
	log.Printf("erro:%s", err)
	return content, false
}

func (r *BookRepo) AddBook(pid string, book *models.Book) bool {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("Error adding book %v\n in repo, can't start transaction", book)
		return false
	}
	_, err = tx.NamedExec("INSERT INTO books (id, title, intro, cover) VALUES (:id, :title, :intro, :cover)", book)
	if err != nil {
		log.Printf("Error adding book %v in repo, with error %v\n", book, err)
		return false
	}
	if book.Authors != nil && len(book.Authors) > 0 {
		for i := range book.Authors {
			author := book.Authors[i]
			aid := author.ID
			_, err = tx.Exec("INSERT INTO books_people (bookid, pid) VALUES ($1, $2)", book.ID, aid)
			if err != nil {
				log.Printf("Error adding book %v in repo, with error %v\n", book, err)
				return false
			}
		}
	} else {
		_, err = tx.Exec("INSERT INTO books_people (bookid, pid) VALUES ($1, $2)", book.ID, pid)
		if err != nil {
			log.Printf("Error adding book %v in repo, with error %v\n", book, err)
			return false
		}
	}
	// add pid to collaborators as initiator
	_, err = tx.Exec("INSERT INTO books_collaborators (bookid, kind, pid) VALUES ($1, $2, $3)", book.ID, 0, pid)
	if err != nil {
		log.Printf("Error adding book %v in repo, with error %v\n", book, err)
		return false
	}

	// add toc into chapters
	// TODO: use batch insert to improve performance
	log.Printf("Got toc:%#v\n", book.Toc)
	toc := book.Toc
	for i := range toc {
		chap := toc[i]
		_, err = tx.Exec("INSERT INTO chapters VALUES ($1, $2, $3, $4)", chap.ID, chap.Order, book.ID, chap.Title)
		if err != nil {
			log.Printf("Error adding chapter %v for book %v in repo, with error %v\n", chap, book.ID, err)
			return false
		}
	}
	tx.Commit()
	return true
}

func (r *BookRepo) ModifyBook(pid string, book *models.Book) bool {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("Error modifying book %v\n in repo, can't start transaction", book)
		return false
	}
	_, err = tx.NamedExec("UPDATE books set title = :title, intro = :intro, cover = :cover where id = :id", book)
	if err != nil {
		log.Printf("Error updating book %v in repo, with error %v\n", book, err)
		return false
	}
	// update the chapters
	err = r.updateChapters(tx, book)
	if err != nil {
		log.Printf("Error updating chapters for book %v in repo, with error %v\n", book, err)
		return false
	}

	/* NOTE: you can't change author of a book, at least for now
	_, err = tx.Exec("INSERT INTO books_people (bookid, pid) VALUES ($1, $2)", book.ID, pid)
	if err != nil {
		log.Printf("Error adding book %v in repo, with error %v\n", book, err)
		return false
	}
	*/
	tx.Commit()
	return true
}

func (r *BookRepo) updateChapters(tx *sqlx.Tx, book *models.Book) error {
	// remove current Chapters
	_, err := tx.Exec("DELETE FROM chapters where bookid = $1", book.ID)
	if err != nil {
		return err
	}
	// add toc into chapters
	// TODO: use batch insert to improve performance
	log.Printf("Got toc:%#v\n", book.Toc)
	toc := book.Toc
	for i := range toc {
		chap := toc[i]
		_, err = tx.Exec("INSERT INTO chapters VALUES ($1, $2, $3, $4)", chap.ID, chap.Order, book.ID, chap.Title)
		if err != nil {
			log.Printf("Error adding chapter %v for book %v in repo, with error %v\n", chap, book.ID, err)
			return err
		}
	}
	return err
}

func (r *BookRepo) ModifyChapter(chapter *models.Chapter) bool {
	// update chapter content
	dir := path.Join(r.rootDir, "books", chapter.BookID[:4], chapter.BookID)
	if err := os.MkdirAll(dir, 0644); err != nil {
		log.Printf("Error creating dir for chapter")
		return false
	}
	file := path.Join(dir, chapter.ID) + ".md"
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Error opening file for chapter")
		return false
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	w.WriteString(chapter.Content)
	w.Flush()

	// update chapter info in db
	{
		_, err := r.db.Exec("UPDATE chapters set title = $1 where bookid = $2 and id = $3", chapter.Title, chapter.BookID, chapter.ID)
		if err != nil {
			log.Printf("Error updating chapter %v with title %v, error: %v\n", chapter.ID, chapter.Title, err)
			return false
		}
	}
	return true
}

func (r *BookRepo) RemoveBook(bookid string) bool {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("Error removing book %v\n in repo, can't start transaction", bookid)
		return false
	}
	_, err = tx.Exec("DELETE from books where id = $1", bookid)
	if err != nil {
		log.Printf("Error removing book %v in repo, with error %v\n", bookid, err)
		return false
	}
	_, err = tx.Exec("DELETE from books_people where bookid = $1", bookid)
	if err != nil {
		log.Printf("Error removing books_people relation %v in repo, with error %v\n", bookid, err)
		return false
	}
	tx.Commit()
	return true
}

func (r *BookRepo) AddBookWithToc(toc *pmodels.Toc) {
	r.db.Exec("INSERT INTO books VALUES($1, $2)", toc.BookID, toc.Title)
	for od, item := range toc.Items {
		_, err := r.db.Exec("INSERT INTO chapters VALUES($1, $2, $3, $4)",
			item.ChapID, od+1, toc.BookID, item.Title)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (r *BookRepo) ListBooksByAuthor(pid string) []models.Book {
	var books []models.Book

	err := r.db.Select(&books, "SELECT b.* from books b, books_people p where b.id = p.bookid and p.pid = $1 order by b.title asc", pid)
	if err != nil {
		log.Printf("Error listing books for pid %v in repo, with error %v\n", pid, err)
		return books
	}

	return books
}

func (r *BookRepo) ListBooksByCollaborator(pid string) []models.Book {
	var books []models.Book

	err := r.db.Select(&books, "SELECT b.* from books b, books_collaborators c where b.id = c.bookid and c.pid = $1 order by b.title asc", pid)
	if err != nil {
		log.Printf("Error listing books for collaborator with pid %v in repo, with error %v\n", pid, err)
		return books
	}

	return books
}

func (r *BookRepo) ListBooksByTranslator(pid string) []models.Book {
	var books []models.Book

	err := r.db.Select(&books, "SELECT b.* from books b, books_collaborators c where b.id = c.bookid and c.pid = $1 and c.kind = 3 order by b.title asc", pid)
	if err != nil {
		log.Printf("Error listing books for collaborator with pid %v in repo, with error %v\n", pid, err)
		return books
	}

	return books
}

func (r *BookRepo) ListPublicBooks(lang string) []models.Book {
	var books []models.Book

	// TODO: add constraint for public books
	var err error
	if lang != "" {
		err = r.db.Select(&books, "SELECT b.* from books b where lang = $1 order by b.title asc", lang)
	} else {
		err = r.db.Select(&books, "SELECT b.* from books b where order by b.title asc")
	}
	if err != nil {
		log.Printf("Error listing public books with lang %v in repo, with error %v\n", lang, err)
		return books
	}

	return books
}

func (r *BookRepo) AddTranslation(bookid string, pid string) bool {
	_, err := r.db.Exec("INSERT INTO books_collaborators VALUES ($1, $2, $3)", bookid, 3, pid)
	if err != nil {
		log.Printf("Error adding translation for book %v and pid %v, with error: %v\n", bookid, pid, err)
		return false
	}
	return true
}
