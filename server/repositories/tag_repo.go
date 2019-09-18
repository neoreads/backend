package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

type TagRepo struct {
	DB *sqlx.DB
}

func NewTagRepo(db *sqlx.DB) *TagRepo {
	return &TagRepo{
		DB: db,
	}
}

var KindMap = map[string]int{
	"topic":   0,
	"event":   1,
	"people":  2,
	"place":   3,
	"time":    4,
	"emotion": 5,
}

func (r *TagRepo) ListNewsTags(t string) []models.Tag {
	tags := []models.Tag{}
	var kind = KindMap[t]

	r.DB.Select(&tags, "SELECT * from tags where kind = $1", kind)
	log.Printf("selected tags: %v\n", tags)
	return tags
}
