package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads-backend/go/server/repositories"
)

// BookController serves book info and book content
type BookController struct {
	Repo *repositories.BookRepo
}

func NewBookController(r *repositories.BookRepo) *BookController {
	return &BookController{Repo: r}
}

func (ctrl *BookController) GetPing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (ctrl *BookController) GetBook(c *gin.Context) {
	bookid := c.Param("bookid")
	book, found := ctrl.Repo.GetBook(bookid)
	if found {
		c.JSON(http.StatusOK, book)
	} else {
		c.String(http.StatusBadRequest, "book info not found!")
	}
}

func (ctrl *BookController) GetContent(c *gin.Context) {
	bookid := c.Param("bookid")
	chapid := c.Param("chapid")
	content, found := ctrl.Repo.GetContent(bookid, chapid)
	if found {
		c.JSON(http.StatusOK, content)
	} else {
		c.String(http.StatusBadRequest, "chapter not found!")
	}
}
