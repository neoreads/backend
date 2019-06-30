package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	app.RegisterView(iris.HTML("./", ".html"))

	// Serve a controller based on the root Router, "/".
	mvc.New(app.Party("/book")).Handle(new(BookController))

	// TODO: read port from config
	app.Run(iris.Addr(":8090"))
}

// BookController serves the "/", "/ping" and "/hello".
type BookController struct{}

// Get serves
// Method:   GET
// Resource: http://localhost:8080
func (c *BookController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

// GetPing serves
// Method:   GET
// Resource: http://localhost:8080/ping
func (c *BookController) GetPing() string {
	return "pong"
}

func (c *BookController) GetPong() string {
	return "ping"
}

// GetHello serves
// Method:   GET
// Resource: http://localhost:8080/hello
func (c *BookController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}

func (c *BookController) GetNow() interface{} {
	return map[string]string{"message": "Now!"}
}

// BeforeActivation called once, before the controller adapted to the main application
// and of course before the server ran.
// After version 9 you can also add custom routes for a specific controller's methods.
// Here you can register custom method's handlers
// use the standard router with `ca.Router` to do something that you can do without mvc as well,
// and add dependencies that will be binded to a controller's fields or method function's input arguments.
func (c *BookController) BeforeActivation(b mvc.BeforeActivation) {
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("Inside /custom_path")
		ctx.Next()
	}
	b.Handle("GET", "/custom_path", "CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)

	// or even add a global middleware based on this controller's router,
	// which in this Book is the root "/":
	// b.Router().Use(myMiddleware)
}

// CustomHandlerWithoutFollowingTheNamingGuide serves
// Method:   GET
// Resource: http://localhost:8080/custom_path
func (c *BookController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}

// GetUserBy serves
// Method:   GET
// Resource: http://localhost:8080/user/{username:string}
// By is a reserved "keyword" to tell the framework that you're going to
// bind path parameters in the function's input arguments, and it also
// helps to have "Get" and "GetBy" in the same controller.
//
func (c *BookController) GetUserBy(username string) mvc.Result {
	return mvc.View{
		Name: "user/username.html",
		Data: username,
	}
}

/* Can use more than one, the factory will make sure
that the correct http methods are being registered for each route
for this controller, uncomment these if you want:
func (c *BookController) Post() {}
func (c *BookController) Put() {}
func (c *BookController) Delete() {}
func (c *BookController) Connect() {}
func (c *BookController) Head() {}
func (c *BookController) Patch() {}
func (c *BookController) Options() {}
func (c *BookController) Trace() {}
*/

/*
func (c *BookController) All() {}
//        OR
func (c *BookController) Any() {}
func (c *BookController) BeforeActivation(b mvc.BeforeActivation) {
	// 1 -> the HTTP Method
	// 2 -> the route's path
	// 3 -> this controller's method name that should be handler for that route.
	b.Handle("GET", "/mypath/{param}", "DoIt", optionalMiddlewareHere...)
}
// After activation, all dependencies are set-ed - so read only access on them
// but still possible to add custom controller or simple standard handlers.
func (c *BookController) AfterActivation(a mvc.AfterActivation) {}
*/
