package storage

import "learn-go/internal/models"

type IngredientStorage interface {
	Create(name string) (*models.Ingredient, error)
	List() ([]*models.Ingredient, error)
	Update(name, newName string) (*models.Ingredient, error)
	SeedTestData() ([]*models.Ingredient, error)
	// TODO: Add GetByID(id int) (*models.Ingredient, error)
	// TODO: Add Delete(id int) error
}
