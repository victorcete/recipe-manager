package models

import (
	"time"
)

type Ingredient struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	KcalPer100g         *float64   `json:"kcal_per_100g,omitempty"`
	FatPer100g          *float64   `json:"fat_per_100g,omitempty"`
	ProteinPer100g      *float64   `json:"protein_per_100g,omitempty"`
	CarbsPer100g        *float64   `json:"carbs_per_100g,omitempty"`
	LastNutritionUpdate *time.Time `json:"last_nutrition_update,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CreateIngredientRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateIngredientRequest struct {
	Name *string `json:"name,omitempty"`
}
