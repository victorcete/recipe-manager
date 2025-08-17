package storage

import "github.com/victorcete/recipe-manager/internal/models"

type IngredientStorage interface {
	Create(name string) (*models.Ingredient, error)
	Delete(name string) error
	List() ([]*models.Ingredient, error)
	SeedTestData() ([]*models.Ingredient, error)
	Update(name, newName string) (*models.Ingredient, error)
}
