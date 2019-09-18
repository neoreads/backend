package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/repositories"
)

type TagController struct {
	Repo *repositories.TagRepo
}

func NewTagController(r *repositories.TagRepo) *TagController {
	return &TagController{
		Repo: r,
	}
}

func (ctrl *TagController) ListNewsTags(c *gin.Context) {
	t := c.Query("t")
	tags := ctrl.Repo.ListNewsTags(t)
	c.JSON(http.StatusOK, tags)
}
