package controllers

import (
	"log"
	"net/http"

	"github.com/neoreads/backend/util"

	"github.com/neoreads/backend/server/models"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/repositories"
)

type UserController struct {
	Repo  *repositories.UserRepo
	IDGen *util.N64Generator
}

func NewUserController(r *repositories.UserRepo) *UserController {
	return &UserController{
		Repo:  r,
		IDGen: util.NewN64Generator(8),
	}
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	var r models.RegisterInfo
	err := c.ShouldBindJSON(&r)
	if err != nil {
		log.Printf("error getting register info:%s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "error getting register info",
		})
		return
	}
	username := r.Username
	// check user exists
	_, found := ctrl.Repo.GetUser(username)
	if found {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "user exists",
		})
		return
	}

	email := r.Email

	// TODO serverside form validation
	log.Printf("registering: u:%s,e:%s\n", username, email)
	// TODO check pid uniqueness in db
	pid := ctrl.IDGen.Next()
	r.Pid = pid

	err = ctrl.Repo.RegisterUser(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    2,
			"message": "error register user in the database",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user register successful",
	})

}
