package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/server/repositories"
)

type PeopleController struct {
	Repo *repositories.PeopleRepo
}

func NewPeopleController(r *repositories.PeopleRepo) *PeopleController {
	return &PeopleController{
		Repo: r,
	}
}

func (ctrl *PeopleController) AddPerson(c *gin.Context) {
	var p models.Person

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	succ := ctrl.Repo.AddPerson(&p)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": p.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}

}
