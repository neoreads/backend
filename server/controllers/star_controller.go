package controllers

import "github.com/neoreads/backend/server/repositories"

type StarController struct {
	Repo *repositories.StarRepo
}

func NewStarController(r *repositories.StarRepo) *StarController {
	return &StarController{
		Repo: r,
	}
}
