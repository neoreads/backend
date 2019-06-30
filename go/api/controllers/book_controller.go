package controllers

import (
	"github.com/kataras/iris/mvc"
	"github.com/neoreads-backend/go/api/datamodels"
	"github.com/neoreads-backend/go/api/services"
)

// BookController serves the "/", "/ping" and "/hello".
type BookController struct {
	Service services.BookService
}

// Get serves
// Method:   GET
// Resource: http://localhost:8080
func (c *BookController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

// Get book by id
// Method: GET
// Resource: /api/book/{id}
func (c *BookController) GetBy(id string) (book datamodels.Book, found bool) {
	/*
		return datamodels.Book{
			ID:     id,
			Title:  "To Kill a Mocking Bird",
			Author: "Harper Lee",
			Desc:   "....",
		}, true
	*/
	var x, succ = c.Service.GetByID(id)
	return x, succ
}

// GetHello serves
// Method:   GET
// Resource: http://localhost:8080/hello
func (c *BookController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}
