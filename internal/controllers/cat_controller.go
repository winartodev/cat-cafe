package controllers

import "github.com/winartodev/cat-cafe/internal/repositories"

type CatController interface {
	GetCatController() (string, error)
}

type catController struct {
	catRepository repositories.CatRepository
}

func NewCatController(catRepo repositories.CatRepository) CatController {
	return &catController{
		catRepository: catRepo,
	}
}

// GetCatController implements CatController.
func (cr *catController) GetCatController() (string, error) {
	return cr.catRepository.GetCatDB()
}
