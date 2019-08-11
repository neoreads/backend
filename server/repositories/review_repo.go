package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

type ReviewRepo struct {
	db *sqlx.DB
}

func NewReviewRepo(db *sqlx.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) ListReviewNotes(pid string, bookid string, chapid string) []models.ReviewNote {
	var notes []models.ReviewNote

	err := r.db.Select(&notes, "SELECT id, ntype, paraid, sentid, content from notes where pid = $1 and bookid = $2 and chapid = $3 and ptype = 1", pid, bookid, chapid)
	if err != nil {
		log.Printf("Error querying for review notes, with error: %v\n", err)
	}

	return notes
}
