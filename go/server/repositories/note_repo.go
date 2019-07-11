package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads-backend/go/server/models"
)

// NoteRepo all kinds of notes
type NoteRepo struct {
	db *sqlx.DB
}

func NewNoteRepo(db *sqlx.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) AddNote(n *models.Note) {
	_, err := r.db.NamedExec("INSERT INTO notes VALUES (:id, :ntype, :ptype, :userid, :bookid, :chapid, :paraid, :sentid, :wordid)", n)
	if err != nil {
		log.Printf("error adding note:%v, with err: %v\n", n, err)
	}
}

func (r *NoteRepo) RemoveNote(id string) {
	_, err := r.db.Exec("DELETE FROM notes where id = $1", id)
	if err != nil {
		log.Printf("error removing note from db:%v, with err:%v\n", id, err)
	}
}

func (r *NoteRepo) ListNotes(userid string, bookid string, chapid string) []models.Note {
	var notes []models.Note
	err := r.db.Select(&notes, "SELECT * from notes where userid = $1 and bookid = $2 and chapid = $3", userid, bookid, chapid)
	if err != nil {
		log.Printf("error listing notes from db:%v, with err:%v\n", userid+"_"+bookid+"_"+chapid, err)
	}
	return notes
}
