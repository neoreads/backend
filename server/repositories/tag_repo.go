package repositories

import (
	"log"
	"strconv"

	sq "github.com/Masterminds/squirrel"

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

func (r *TagRepo) ListTags(role string, kind string) []models.Tag {
	tags := []models.Tag{}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select("*").From("tags")

	if role != "" {
		c, err := strconv.Atoi(role)
		if err != nil {
			log.Printf("Error converting parameter role %v to int\n", role)
			return tags
		}
		q = q.Where("role = ?", c)
	}
	if kind != "" {
		k, err := strconv.Atoi(kind)
		if err != nil {
			log.Printf("Error converting parameter kind %v to int\n", kind)
			return tags
		}
		q = q.Where("kind = ?", k)
	}

	sql, args, err := q.ToSql()
	log.Printf("Got sql: %v, %#v\n", sql, args)
	if err != nil {
		log.Printf("err building sql for ListTags : %v\n", err)
	}

	err = r.DB.Select(&tags, sql, args...)
	if err != nil {
		log.Printf("error selecting from tags: %v, with error: %v\n", sql, err)
		return tags
	}

	log.Printf("selected tags: %v\n", tags)
	return tags
}
