package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/repositories"
)

type ReviewController struct {
	Repo *repositories.ReviewRepo
}

func NewReviewController(r *repositories.ReviewRepo) *ReviewController {
	return &ReviewController{Repo: r}
}

func (ctrl *ReviewController) ListReviewNotes(c *gin.Context) {
	pid := "00000001"
	bookid := c.Param("bookid")
	chapid := c.Param("chapid")
	notes := ctrl.Repo.ListReviewNotes(pid, bookid, chapid)
	c.JSON(http.StatusOK, notes)
}
