package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Student a data type for students
type Student struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=hello sslmode=disable password=123456")
	if err != nil {
		log.Fatal(err)
	}

	students := []Student{}
	db.Select(&students, "SELECT * FROM student")

	for s := range students {
		fmt.Println(students[s])
	}
}
