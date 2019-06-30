package main

import (
	"strconv"

	"github.com/kataras/iris"

	"github.com/kataras/iris/hero"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func hello(to string) string {
	return "Hello: " + to
}

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
	})

	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	app.Get("/hello.json", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	app.Get("/users/{id:uint64}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		//ctx.HTML("<h1>" + strconv.FormatUint(id, 10) + "</h1>")
		ctx.HTML("<h1>" + strconv.FormatUint(id, 10) + "</h1>")
	})

	app.Get("/welcome", func(ctx iris.Context) {
		name := ctx.URLParamDefault("name", "visus")
		ctx.HTML("Hello: " + name)
	})

	app.Get("/hello/{to:string}", hero.Handler(hello))

	app.Run(iris.Addr(":8089"))
}
