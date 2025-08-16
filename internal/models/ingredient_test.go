package models

import (
	"testing"
	"time"
)

func TestNewIngredient(t *testing.T) {
	t.Run("successful model creation", func(t *testing.T) {
		id := 1
		name := "tomato"

		ingredient := NewIngredient(id, name)

		if ingredient.ID != id {
			t.Errorf("expected ID %d, got %d", id, ingredient.ID)
		}

		if ingredient.Name != name {
			t.Errorf("expected name %q, got %q", name, ingredient.Name)
		}

		if time.Since(ingredient.CreatedAt) > time.Second {
			t.Errorf("expected creation date to be recent, got %v", ingredient.CreatedAt)
		}
	})

	t.Run("creates different ingredients", func(t *testing.T) {
		ingredient1 := NewIngredient(1, "tomato")
		ingredient2 := NewIngredient(2, "potato")
		if ingredient1.ID == ingredient2.ID {
			t.Errorf("expected different IDs for different ingredients")
		}
		if ingredient1.Name == ingredient2.Name {
			t.Errorf("expected different names for different ingredients")
		}
	})
}
