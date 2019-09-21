package controllers

import (
	"net/http"
	"strings"

	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/util"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/repositories"
)

// BookController serves book info and book content
type BookController struct {
	Repo      *repositories.BookRepo
	IDGen     *util.N64Generator
	ChapIDGen *util.N64Generator
	ParaIDGen *util.N64Generator
	SentIDGen *util.N64Generator
}

func NewBookController(r *repositories.BookRepo) *BookController {
	return &BookController{
		Repo:      r,
		IDGen:     util.NewN64Generator(8),
		ChapIDGen: util.NewN64Generator(4),
		ParaIDGen: util.NewN64Generator(4),
		SentIDGen: util.NewN64Generator(4),
	}
}

func (ctrl *BookController) GetBook(c *gin.Context) {
	// TODO: add authorization, because not all books are public
	bookid := c.Param("bookid")
	book, found := ctrl.Repo.GetBook(bookid)
	if found {
		c.JSON(http.StatusOK, book)
	} else {
		c.String(http.StatusBadRequest, "book info not found!")
	}
}

func (ctrl *BookController) RemoveBook(c *gin.Context) {
	// TODO: add authorization, because not all books are public
	bookid := c.Param("bookid")
	succ := ctrl.Repo.RemoveBook(bookid)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book info not found!"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "chapter not found!"})
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

	// gen chap ids
	for i := range book.Toc {
		chap := &book.Toc[i]
		if chap.ID == "" {
			chap.ID = ctrl.ChapIDGen.Next()
		}
	}

	succ := ctrl.Repo.AddBook(pid, &book)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": bookid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}

}

func (ctrl *BookController) ModifyBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid

	// gen chap ids
	for i := range book.Toc {
		chap := &book.Toc[i]
		if strings.TrimSpace(chap.ID) == "" {
			chap.ID = ctrl.ChapIDGen.Next()
		}
	}

	succ := ctrl.Repo.ModifyBook(pid, &book)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": book.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}
}

func (ctrl *BookController) parseChapterContent(content string) string {
	paras := util.ParseMD(content)
	sb := strings.Builder{}

	for i := range paras {
		p := paras[i]
		sents := p.Sents
		for j := range sents {
			sent := sents[j]
			sb.WriteString(sent.Content)
			sb.WriteString("{s:" + sent.ID + "}")
		}
		sb.WriteString("{p:" + p.ID + "}")
		sb.WriteString("\n")
		sb.WriteString("\n")
	}
	return sb.String()
}

func (ctrl *BookController) ModifyChapter(c *gin.Context) {
	var chapter models.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: check if the user is authorized to do this action
	// user, _ := c.Get("jwtuser")
	// pid := user.(*models.User).Pid

	// filter chapter content to add para/sent ids

	content := chapter.Content
	chapter.Content = ctrl.parseChapterContent(content)

	succ := ctrl.Repo.ModifyChapter(&chapter)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": chapter.ID, "bookid": chapter.BookID})
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

func (ctrl *BookController) ListMyCollaborationBooks(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid

	books := ctrl.Repo.ListBooksByCollaborator(pid)
	c.JSON(http.StatusOK, books)
}

func (ctrl *BookController) ListPublicBooks(c *gin.Context) {
	lang := c.Query("lang")
	books := ctrl.Repo.ListPublicBooks(lang)
	c.JSON(http.StatusOK, books)
}

func (ctrl *BookController) AddTranslation(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid
	bookid := c.Param("bookid")

	succ := ctrl.Repo.AddTranslation(bookid, pid)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "bookid": bookid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}
}

func (ctrl *BookController) ListMyTranslationBooks(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid

	books := ctrl.Repo.ListBooksByTranslator(pid)
	c.JSON(http.StatusOK, books)
}
