package controllers

import (
	"net/http"

	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/util"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/repositories"
)

// BookController serves book info and book content
type BookController struct {
	Repo  *repositories.BookRepo
	IDGen *util.N64Generator
}

func NewBookController(r *repositories.BookRepo) *BookController {
	return &BookController{
		Repo:  r,
		IDGen: util.NewN64Generator(8),
	}
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

func (ctrl *BookController) GetTOC(c *gin.Context) {
	bookid := c.Param("bookid")
	toc := ctrl.Repo.GetTOC(bookid)

	c.JSON(http.StatusOK, toc)
}

func (ctrl *BookController) GetBookChapter(c *gin.Context) {
	bookid := c.Param("bookid")
	chapid := c.Param("chapid")
	content, found := ctrl.Repo.GetContent(bookid, chapid)
	if found {
		c.JSON(http.StatusOK, content)
	} else {
		c.String(http.StatusBadRequest, "chapter not found!")
	}
}

func (ctrl *BookController) HotList(c *gin.Context) {
	var books []models.Book
	books = append(books, models.Book{
		ID:    "bKbnk8Zd",
		Title: "史记",
	})
	c.JSON(http.StatusOK, books)
}

func (ctrl *BookController) AddBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// gen book id
	bookid := ctrl.IDGen.Next()
	book.ID = bookid

	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid

	succ := ctrl.Repo.AddBook(pid, &book)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": bookid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}

}

func (ctrl *BookController) ListMyBooks(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid

	books := ctrl.Repo.ListBooksByAuthor(pid)
	c.JSON(http.StatusOK, books)
}
