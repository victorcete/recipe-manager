package models

import "time"

// Ingredient represents a cooking ingredient.
type Ingredient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewIngredient creates a new ingredient.
func NewIngredient(id int, name string) *Ingredient {
	now := time.Now()
	return &Ingredient{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
