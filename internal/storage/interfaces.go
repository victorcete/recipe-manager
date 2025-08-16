package storage

import "learn-go/internal/models"

type IngredientStorage interface {
	Create(name string) (*models.Ingredient, error)
}
