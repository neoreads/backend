package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neoreads-backend/go/server/repositories"
)

type UserController struct {
	Repo *repositories.UserRepo
}

func NewUserController(r *repositories.UserRepo) *UserController {
	return &UserController{
		Repo: r,
	}
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	// check user exists
	_, found := ctrl.Repo.GetUser(username)
	if found {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "user exists",
		})
		return
	}

	password := c.PostForm("password")
	email := c.PostForm("email")
	err := ctrl.Repo.RegisterUser(username, email, password)
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
