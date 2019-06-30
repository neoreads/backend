package main

import (
	"github.com/kataras/iris"
	"github.com/neoreads-backend/go/api/controllers"
	"github.com/neoreads-backend/go/api/services"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	mvc.Configure(app.Party("/api/book"), books)

	// TODO: read port from config
	app.Run(iris.Addr(":8090"))
}

func books(app *mvc.Application) {
	bookService := services.NewBookService()
	app.Register(bookService)

	app.Handle(new(controllers.BookController))
}
