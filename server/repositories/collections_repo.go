package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

// CollectionRepo Collection related data repository
type CollectionRepo struct {
	db *sqlx.DB
}

// NewCollectionRepo creator for CollectionRepo
func NewCollectionRepo(db *sqlx.DB) *CollectionRepo {
	return &CollectionRepo{db: db}
}

func (r *CollectionRepo) AddCollection(col *models.Collection) bool {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("error adding Collection:%v, with err: %v\n", col, err)
		return false
	}
	_, err = tx.NamedExec("INSERT INTO collections (id, title, intro)"+
		" VALUES (:id, :title, :intro)", col)
	if err != nil {
		log.Printf("error adding Collection:%v, with err: %v\n", col, err)
		return false
	}
	_, err = tx.Exec("INSERT INTO collections_people (colid, pid) VALUES ($1, $2)", col.ID, col.PID)
	if err != nil {
		log.Printf("error adding Collection:%v, with err: %v\n", col, err)
		return false
	}

	artids := col.ArtIDs
	if len(artids) > 0 {
		for i := range artids {
			artid := artids[i]
			_, err = tx.Exec("INSERT INTO collections_articles (colid, artid) VALUES ($1, $2)", col.ID, artid)
			if err != nil {
				log.Printf("error adding Collection:%v, with err: %v\n", col, err)
				return false
			}
		}
	}
	tx.Commit()
	return true
}

func (r *CollectionRepo) ModifyCollection(c *models.Collection) bool {
	// TODO: add support for modTime
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("error modifying Collection:%v, with err: %v\n", c, err)
		return false
	}
	_, err = tx.NamedExec("UPDATE Collections set title = :title, intro = :intro, modtime = now() where id = :id", c)
	if err != nil {
		log.Printf("error modifying Collection:%v, with err: %v\n", c, err)
		return false
	}
	// update collections_articles
	_, err = tx.Exec("DELETE FROM collections_articles where colid = $1", c.ID)

	artids := c.ArtIDs
	if len(artids) > 0 {
		for i := range artids {
			artid := artids[i]
			_, err = r.db.Exec("INSERT INTO collections_articles (colid, artid) VALUES ($1, $2)", c.ID, artid)
			if err != nil {
				log.Printf("error adding Collection:%v, with err: %v\n", c, err)
				return false
			}
		}
	}
	tx.Commit()
	return true
}

func (r *CollectionRepo) GetCollection(colid string) models.Collection {
	var collection models.Collection
	err := r.db.Get(&collection, "SELECT c.title, c.intro, c.id, p.pid, array_remove(array_agg(a.artid), NULL)::text[] as artids from collections c join collections_people p on c.id = p.colid and c.id = $1 left join collections_articles a on c.id = a.colid group by c.id, p.pid order by c.modtime desc", colid)
	if err != nil {
		log.Printf("error listing collections from db:%v, with err:%v\n", colid, err)
	}
	return collection
}

func (r *CollectionRepo) ListCollections(pid string) []models.Collection {
	collections := []models.Collection{}
	// TODO: this might be a performance nightmare
	sql := "SELECT c.*, p.pid, array_remove(array_agg(a.artid), NULL)::text[] as artids from collections c join collections_people p on c.id = p.colid and p.pid = $1" + " left join collections_articles a on c.id = a.colid group by c.id, p.pid order by c.modtime desc;"
	err := r.db.Select(&collections, sql, pid)
	for a := range collections {
		col := collections[a]
		log.Printf("col:%v\n", col)
	}
	if err != nil {
		log.Printf("error listing Collections from db:%v, with err:%v\n", pid, err)
	}
	return collections
}

func (r *CollectionRepo) RemoveCollection(colid string) bool {
	tx := r.db.MustBegin()
	_, err := tx.Exec("DELETE FROM collections where id = $1", colid)
	if err != nil {
		log.Printf("error removing collection from db:%v, with err:%v\n", colid, err)
		return false
	}
	_, err = tx.Exec("DELETE FROM collections_people where colid = $1", colid)
	if err != nil {
		log.Printf("error removing collections_people relation from db:%v, with err:%v\n", colid, err)
		return false
	}

	_, err = tx.Exec("DELETE FROM collections_articles where colid = $1", colid)
	if err != nil {
		log.Printf("error removing collections_articles relation from db:%v, with err:%v\n", colid, err)
		return false
	}
	tx.Commit()
	return true
}
