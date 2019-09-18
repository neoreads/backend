package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/server/repositories"
	"github.com/neoreads/backend/util"
)

type NewsController struct {
	Repo  *repositories.NewsRepo
	IDGen *util.N64Generator
}

func NewNewsController(r *repositories.NewsRepo) *NewsController {
	return &NewsController{
		Repo:  r,
		IDGen: util.NewN64Generator(8),
	}
}

func (ctrl *NewsController) AddNews(c *gin.Context) {
	var post models.News

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// gen news id
	newsid := ctrl.IDGen.Next()
	post.ID = newsid

	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid
	post.PID = pid

	succ := ctrl.Repo.AddNews(&post)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": newsid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}
}

func (ctrl *NewsController) ListNews(c *gin.Context) {
	tagid := c.Query("tagid")
	newsList := ctrl.Repo.ListNews(tagid)
	c.JSON(http.StatusOK, newsList)
}
