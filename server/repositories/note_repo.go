package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

// NoteRepo all kinds of notes
type NoteRepo struct {
	db *sqlx.DB
}

func NewNoteRepo(db *sqlx.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) AddNote(n *models.Note) {
	_, err := r.db.NamedExec("INSERT INTO notes (id, ntype, ptype, pid, bookid, chapid, paraid, sentid, wordid, content)"+
		" VALUES (:id, :ntype, :ptype, :pid, :bookid, :chapid, :paraid, :sentid, :wordid, :content)", n)
	if err != nil {
		log.Printf("error adding note:%v, with err: %v\n", n, err)
	}
}

func (r *NoteRepo) ModifyNote(n *models.Note) bool {
	_, err := r.db.NamedExec("UPDATE notes set content = :content where id = :id", n)
	if err != nil {
		log.Printf("error modifying note:%v, with err: %v\n", n, err)
		return false
	}
	return true
}

func (r *NoteRepo) RemoveNote(id string) {
	_, err := r.db.Exec("DELETE FROM notes where id = $1", id)
	if err != nil {
		log.Printf("error removing note from db:%v, with err:%v\n", id, err)
	}
}

func (r *NoteRepo) ListNotes(pid string, bookid string, chapid string) []models.Note {
	notes := []models.Note{}
	err := r.db.Select(&notes, "SELECT * from notes where pid = $1 and bookid = $2 and chapid = $3", pid, bookid, chapid)
	if err != nil {
		log.Printf("error listing notes from db:%v, with err:%v\n", pid+"_"+bookid+"_"+chapid, err)
	}
	return notes
}
