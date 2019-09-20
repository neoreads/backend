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

func (ctrl *TagController) ListTags(c *gin.Context) {
	class := c.Query("c")
	kind := c.Query("k")

	tags := ctrl.Repo.ListTags(class, kind)
	c.JSON(http.StatusOK, tags)
}
