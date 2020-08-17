package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/util"
)

// NoteRepo all kinds of notes
type NoteRepo struct {
	db        *sqlx.DB
	LineIDGen *util.N64Generator
}

func NewNoteRepo(db *sqlx.DB) *NoteRepo {
	return &NoteRepo{db: db, LineIDGen: util.NewN64Generator(4)}
}

func (r *NoteRepo) AddNote(n *models.Note) bool {
	// 如果是新的位置，需要给当前行添加行ID。
	log.Printf("lineID for current note: %s\n", n.LineID)
	if n.LineID == "" {
		log.Printf("generating new line id")
		r.GenLineID(n)
	}

	// 如果是现成位置，直接利用现行ID
	_, err := r.db.NamedExec("INSERT INTO notes (id, ntype, ptype, pid, colid, artid, lineid, startpos, endpos, content, value)"+
		" VALUES (:id, :ntype, :ptype, :pid, :colid, :artid, :lineid, :startpos, :endpos, :content, :value)", n)
	if err != nil {
		log.Printf("error adding note:%v, with err: %v\n", n, err)
		return false
	}
	return true
}

func (r *NoteRepo) GenLineID(n *models.Note) bool {
	lid := r.LineIDGen.Next()
	n.LineID = lid

	_, err := r.db.Exec("INSERT INTO lineid_linenum (colid, artid, linenum, lineid) VALUES ($1, $2, $3, $4)", n.ColID, n.ArtID, n.LineNum, n.LineID)
	if err != nil {
		log.Printf("error generating lineid:%v, with err: %v\n", n, err)
		return false
	}
	return true
}

func (r *NoteRepo) GetStars(noteid string) int {
	var value int
	err := r.db.Get(&value, "SELECT value from notes where id = $1", noteid)
	if err != nil {
		log.Printf("Error getting stars for %v, with err: %v\n", noteid, err)
		return 0
	}
	return value
}

func (r *NoteRepo) UpdateStars(noteid string, delta int) bool {
	_, err := r.db.Exec("UPDATE notes set value = value + $1 where id = $2", delta, noteid)
	if err != nil {
		log.Printf("error updating star:%v, with err: %v\n", noteid, err)
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

func (r *NoteRepo) ListStarsForArticles(pid string, artids []string) []models.Note {
	notes := []models.Note{}

	query, args, err := sqlx.In("SELECT * from notes where pid = ? and artid in (?) and ntype = 0 and ptype = 3", pid, artids)
	//query, args, err := sqlx.In("SELECT * from notes where artid in (?) and ntype = 0 and ptype = 3", artids)
	if err != nil {
		log.Printf("error making in query from db:%v[%#v], with err:%v\n", pid, artids, err)
	}
	query = r.db.Rebind(query)
	err = r.db.Select(&notes, query, args...)
	if err != nil {
		log.Printf("error listing notes from db:%v[%#v], with err:%v\n", pid, artids, err)
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
