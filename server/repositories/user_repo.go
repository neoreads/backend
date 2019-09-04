package repositories

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

// UserRepo User related data repository
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo creator for UserRepo
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetUser get user with username
func (r *UserRepo) GetUser(username string) (user models.User, found bool) {
	err := r.db.Get(&user, "SELECT id, pid, username from users where username = $1", username)
	if err != nil {
		log.Printf("err getting userinfo %s, with error: %s\n", username, err)
		return models.User{}, false
	}
	return user, true
	/*
		err := r.db.Get(&user, "SELECT username, firstname, lastname from users_people where username = $1", username)
		if err != nil {
			log.Printf("err getting userinfo %s, with error: %s\n", username, err)
			return models.User{}, false
		}
		return user, true
	*/
}

// CheckLogin check if a user has the right password
// TODO: check password hashes instead of plain password
func (r *UserRepo) CheckLogin(username string, password string) bool {
	var pwd string
	log.Printf("username: %s, password:%s\n", username, password)
	err := r.db.Get(&pwd, "SELECT pwd from users where username = $1", username)
	log.Printf("pwd:%s", pwd)
	if err != nil {
		log.Printf("error getting user %s, with err: %s\n", username, err)
		return false
	}
	log.Printf("=? :%v\n", pwd == strings.TrimSpace(password))
	return pwd == strings.TrimSpace(password)
}

func (r *UserRepo) RegisterUser(reg *models.RegisterInfo) error {
	// create a record in people table
	sql := "INSERT INTO people (id, firstname, lastname) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(sql, reg.Pid, reg.FirstName, reg.LastName)
	if err != nil {
		log.Printf("error registering person %v, with err:%s\n", r, err)
		return err
	}
	// create a record in users table
	sql = "INSERT INTO users (username, email, pwd, pid) VALUES ($1, $2, $3, $4)"
	_, err = r.db.Exec(sql, reg.Username, reg.Email, reg.Password, reg.Pid)
	if err != nil {
		log.Printf("error registering user %s, with err:%s\n", reg.Username, err)
		return err
	}

	return nil
}
