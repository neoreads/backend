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

func (r *NoteRepo) AddNote(n *models.Note) bool {
	_, err := r.db.NamedExec("INSERT INTO notes (id, ntype, ptype, pid, colid, artid, paraid, sentid, startpos, endpos, content)"+
		" VALUES (:id, :ntype, :ptype, :pid, :colid, :artid, :paraid, :sentid, :startpos, :endpos, :content)", n)
	if err != nil {
		log.Printf("error adding note:%v, with err: %v\n", n, err)
		return false
	}
	return true
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

func (r *NoteRepo) ListNotesForPid(pid string, colid string, artid string) []models.Note {
	notes := []models.Note{}
	err := r.db.Select(&notes, "SELECT * from notes where pid = $1 and colid = $2 and artid = $3", pid, colid, artid)
	if err != nil {
		log.Printf("error listing notes from db:%v, with err:%v\n", pid+"_"+colid+"_"+artid, err)
	}
	return notes
}

func (r *NoteRepo) ListNotes(colid string, artid string) []models.Note {
	notes := []models.Note{}
	err := r.db.Select(&notes, "SELECT * from notes where colid = $1 and artid = $2", colid, artid)
	if err != nil {
		log.Printf("error listing notes from db:%v, with err:%v\n", colid+":"+artid, err)
	}
	return notes
}

func (r *NoteRepo) ListNoteCards(colid string, artid string) []models.NoteCard {
	cards := []models.NoteCard{}
	err := r.db.Select(&cards, "SELECT n.*, p.fullname as pname from notes n, people p where n.pid = p.id and n.colid = $1 and n.artid = $2", colid, artid)
	if err != nil {
		log.Printf("error listing notes from db:%v, with err:%v\n", colid+":"+artid, err)
	}
	return cards
}
