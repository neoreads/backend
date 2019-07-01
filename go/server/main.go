package main

import (
	"log"

	"github.com/neoreads-backend/go/server/controllers"

	"github.com/neoreads-backend/go/server/repositories"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var db *sqlx.DB

func initDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=neoreads sslmode=disable password=123456")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	api := r.Group("/api")
	// BOOK
	setupBook(api)
	return r
}

func setupBook(rg *gin.RouterGroup) {

	repo := repositories.NewBookRepo(db)
	ctrl := controllers.NewBookController(repo)

	rg.GET("/ping", ctrl.GetPing)
	rg.GET("/book/:bookid", ctrl.GetBook)
	rg.GET("/book/:bookid/:chapid", ctrl.GetContent)
}

func main() {

	db = initDB()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8090")
}
