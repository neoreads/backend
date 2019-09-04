package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/server/repositories"
	"github.com/neoreads/backend/util"
)

type ArticleController struct {
	Repo  *repositories.ArticleRepo
	IDGen *util.N64Generator
}

func NewArticleController(r *repositories.ArticleRepo) *ArticleController {
	return &ArticleController{
		Repo:  r,
		IDGen: util.NewN64Generator(8),
	}
}

func (ctrl *ArticleController) ModifyArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: check if PID from credential is the same as claimed in the post data
	succ := ctrl.Repo.ModifyArticle(&article)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": article.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error adding article in repo"})
	}
}
func (ctrl *ArticleController) AddArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aid := ctrl.IDGen.Next()
	article.ID = aid
	user, _ := c.Get("jwtuser")
	article.PID = user.(*models.User).Pid

	succ := ctrl.Repo.AddArticle(&article)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": aid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error adding article in repo"})
	}
}

func (ctrl *ArticleController) GetArticle(c *gin.Context) {
	artid := c.Param("artid")
	article := ctrl.Repo.GetArticle(artid)
	c.JSON(http.StatusOK, article)
}

func (ctrl *ArticleController) ListArticles(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid
	articles := ctrl.Repo.ListArticles(pid)
	c.JSON(http.StatusOK, articles)
}

func (ctrl *ArticleController) ListArticlesInCollection(c *gin.Context) {
	colid := c.Param("colid")
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid
	articles := ctrl.Repo.ListArticlesInCollection(pid, colid)
	c.JSON(http.StatusOK, articles)
}

func (ctrl *ArticleController) RemoveArticle(c *gin.Context) {
	id := c.Param("artid")
	succ := ctrl.Repo.RemoveArticle(id)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error removing article in repo"})
	}
}
