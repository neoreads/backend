package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

type PeopleRepo struct {
	DB *sqlx.DB
}

func NewPeopleRepo(db *sqlx.DB) *PeopleRepo {
	return &PeopleRepo{
		DB: db,
	}
}

func (r *PeopleRepo) AddPerson(person *models.Person) bool {

	// TODO: 检查是否有名字冲突
	var id string
	rs, err := r.DB.NamedQuery("INSERT INTO people (fullname, othernames, intro, avatar) VALUES (:fullname, :othernames, :intro, :avatar) RETURNING id", &person)
	if err != nil {
		log.Printf("Error inserting person: %#v, with error: %v\n", person, err)
		return false
	}
	rs.Next()
	rs.Scan(&id)
	log.Printf("person added:%v\n", id)
	person.ID = id
	return true
}

func (r *PeopleRepo) ModifyPerson(personForm *models.PersonForm) bool {
	tx, err := r.DB.Beginx()
	log.Printf("Got forms: %#v,\n", personForm.Tags)
	_, err = tx.NamedExec("UPDATE people set fullname = :fullname, othernames = :othernames, intro = :intro, avatar = :avatar where id = :id", &personForm)
	if err != nil {
		log.Printf("Error modifying person: %#v, with error: %v\n", personForm, err)
		return false
	}
	tx.Exec("DELETE from people_tags where pid = $1", personForm.ID)
	for i := range personForm.Tags {
		tag := personForm.Tags[i]
		tx.Exec("INSERT INTO people_tags VALUES ($1, $2)", personForm.ID, tag.ID)
	}
	log.Printf("person modified:%v\n", personForm.ID)
	tx.Commit()
	return true
}

func (r *PeopleRepo) ListPeople() (people []models.Person) {
	err := r.DB.Select(&people, "SELECT id, fullname, othernames, intro, avatar from people order by fullname")
	if err != nil {
		log.Printf("Error listing people, with error:%v\n", err)
	}
	return people
}

func (r *PeopleRepo) GetPerson(pid string) (person models.PersonForm) {
	err := r.DB.Get(&person, "SELECT id, fullname, othernames, intro, avatar from people where id = $1", pid)
	var tags []models.Tag
	r.DB.Select(&tags, "SELECT t.* from tags t, people_tags pt where t.id = pt.tid and pt.pid = $1", pid)
	if len(tags) > 0 {
		person.Tags = tags
	}
	if err != nil {
		log.Printf("Error listing people, with error:%v\n", err)
	}
	return person
}

func (r *PeopleRepo) HotAuthors(tag string) (people []models.Person) {
	var tagid string
	err := r.DB.Get(&tagid, "SELECT id from tags where role = 3 and tag = $1", tag)
	if err != nil {
		log.Printf("Error finding tag %v, with error:%v\n", tag, err)
		return people
	}
	log.Printf("found tag id %v for tag name: %v", tagid, tag)
	err = r.DB.Select(&people, "SELECT id, fullname, othernames, intro, avatar from people p, people_tags pt where pt.pid = p.id and pt.tid = $1", tagid)
	if err != nil {
		log.Printf("Error listing people, with error:%v\n", err)
	}
	return people
}
