package models

import "time"

// Ingredient represents a cooking ingredient.
// TODO: Add calories_per_100g field for nutritional tracking
// TODO: Future - add protein, fat, carbs per 100g
// TODO: Future - add unit conversion fields (piece_weight_grams, piece_name, density_g_per_ml)
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
