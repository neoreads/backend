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

func (ctrl *PeopleController) ModifyPerson(c *gin.Context) {
	var pf models.PersonForm

	if err := c.ShouldBindJSON(&pf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	succ := ctrl.Repo.ModifyPerson(&pf)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": pf.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo error"})
	}

}

func (ctrl *PeopleController) ListPeople(c *gin.Context) {
	people := ctrl.Repo.ListPeople()
	c.JSON(http.StatusOK, people)
}

func (ctrl *PeopleController) GetPerson(c *gin.Context) {
	pid := c.Param("pid")
	person := ctrl.Repo.GetPerson(pid)
	c.JSON(http.StatusOK, person)
}

func (ctrl *PeopleController) HotAuthors(c *gin.Context) {
	tag := c.Query("tag")
	authors := ctrl.Repo.HotAuthors(tag)
	c.JSON(http.StatusOK, authors)
}
