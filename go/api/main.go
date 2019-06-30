package main

import (
	"log"

	"github.com/neoreads-backend/go/api/repositories"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris"
	_ "github.com/lib/pq"
	"github.com/neoreads-backend/go/api/controllers"
	"github.com/neoreads-backend/go/api/services"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
)

func initDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=neoreads sslmode=disable password=123456")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var database *sqlx.DB

func main() {
	database = initDB()
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	mvc.Configure(app.Party("/api/book"), books)

	// TODO: read port from config
	app.Run(iris.Addr(":8090"))
}

func books(app *mvc.Application) {
	bookRepo := repositories.NewBookRepo(database)
	bookService := services.NewBookService(bookRepo)
	app.Register(bookService)

	app.Handle(new(controllers.BookController))
}
