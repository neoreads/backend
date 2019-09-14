package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

// NewsRepo
type NewsRepo struct {
	db *sqlx.DB
}

// NewsRepo creator
func NewNewsRepo(db *sqlx.DB) *NewsRepo {
	return &NewsRepo{db: db}
}

func (r *NewsRepo) AddNews(news *models.News) bool {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("error adding news item, can't create transaction: %v\n ", err)
		return false
	}
	_, err = tx.NamedExec("INSERT INTO news (id, kind, link, source, title, summary, content) VALUES (:id, :kind, :link, :source, :title, :summary, :content)", news)
	if err != nil {
		log.Printf("error adding news item, insert error: %v\n ", err)
		return false
	}
	tx.Commit()
	return true
}

func (r *NewsRepo) ListNews() []models.News {
	var news []models.News
	err := r.db.Select(&news, "SELECT * from news")
	if err != nil {
		log.Printf("Error listing news: %v\n", err)
	}
	return news
}
