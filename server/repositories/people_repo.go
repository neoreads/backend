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

	// TODO: check name conflicts
	var id string
	rs, err := r.DB.NamedQuery("INSERT INTO people (fullname, othernames, intro, avatar) VALUES (:fullname, :othernames, :intro, :avatar) RETURNING id", &person)
	if err != nil {
		log.Printf("Error inserting person: %#v, with error: %v\n", person, err)
		return false
	}
	rs.Scan(&id)
	person.ID = id
	return true
}
