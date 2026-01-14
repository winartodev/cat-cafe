package controllers

import "github.com/winartodev/cat-cafe/internal/repositories"

type Controller struct {
	CatController CatController
}

func SetUpController(repo repositories.Repository) *Controller {

	catController := NewCatController(repo.CatRepository)

	return &Controller{
		CatController: catController,
	}
}
