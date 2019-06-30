package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type student struct {
	id   int64
	name string
}

func main() {
	connStr := "user=postgres dbname=hello sslmode=disable password=123456"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var s student
		err = rows.Scan(&s.id, &s.name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%v]:%v\n", s.id, s.name)
	}
}
