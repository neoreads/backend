package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/server/repositories"
	"github.com/neoreads/backend/util"
)

type CollectionController struct {
	Repo  *repositories.CollectionRepo
	IDGen *util.N64Generator
}

func NewCollectionController(r *repositories.CollectionRepo) *CollectionController {
	return &CollectionController{
		Repo:  r,
		IDGen: util.NewN64Generator(8),
	}
}

func (ctrl *CollectionController) ModifyCollection(c *gin.Context) {
	var collection models.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: check if PID from credential is the same as claimed in the post data
	succ := ctrl.Repo.ModifyCollection(&collection)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": collection.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error adding Collection in repo"})
	}
}

func (ctrl *CollectionController) AddCollection(c *gin.Context) {
	var Collection models.Collection
	if err := c.ShouldBindJSON(&Collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aid := ctrl.IDGen.Next()
	Collection.ID = aid
	user, _ := c.Get("id")
	Collection.PID = user.(*models.Credential).Username

	succ := ctrl.Repo.AddCollection(&Collection)

	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": aid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error adding Collection in repo"})
	}
}

func (ctrl *CollectionController) GetCollection(c *gin.Context) {
	colid := c.Param("colid")
	collection := ctrl.Repo.GetCollection(colid)
	c.JSON(http.StatusOK, collection)
}

func (ctrl *CollectionController) ListCollections(c *gin.Context) {
	user, _ := c.Get("id")
	username := user.(*models.Credential).Username
	Collections := ctrl.Repo.ListCollections(username)
	c.JSON(http.StatusOK, Collections)
}

func (ctrl *CollectionController) RemoveCollection(c *gin.Context) {
	id := c.Param("colid")
	succ := ctrl.Repo.RemoveCollection(id)
	if succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error removing Collection in repo"})
	}
}
