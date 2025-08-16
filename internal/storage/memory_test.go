package storage

import (
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()
	if storage == nil {
		t.Fatal("memory storage created is nil")
	}

	if storage.ingredients == nil {
		t.Error("ingredients map is nil")
	}

	if storage.nextID != 1 {
		t.Errorf("expected ID to be 1, got %d", storage.nextID)
	}
}

func TestCreateIngredient(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		storage := NewMemoryStorage()
		ingredient, err := storage.Create("tomato")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ingredient.ID != 1 {
			t.Errorf("expected ID 1, got %d", ingredient.ID)
		}

		if ingredient.Name != "tomato" {
			t.Errorf("expected name 'tomato', got '%s'", ingredient.Name)
		}

		if ingredient.CreatedAt.IsZero() {
			t.Errorf("creation date should not be zero")
		}

		if ingredient.UpdatedAt.IsZero() {
			t.Errorf("update date should not be zero")
		}
	})

	t.Run("duplicated ingredient names", func(t *testing.T) {
		storage := NewMemoryStorage()
		ingredient1 := "tomato"
		ingredient2 := "tomato"

		_, err := storage.Create(ingredient1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = storage.Create(ingredient2)
		if err != ErrIngredientNameExists {
			t.Errorf("expected %v, got %v", ErrIngredientNameExists, err)
		}
	})

	t.Run("case insensitive duplicates", func(t *testing.T) {
		storage := NewMemoryStorage()
		ingredient1 := "tomato"
		ingredient2 := " TOMATO "

		_, err := storage.Create(ingredient1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = storage.Create(ingredient2)
		if err != ErrIngredientNameExists {
			t.Errorf("expected %v, got %v", ErrIngredientNameExists, err)
		}
	})

	t.Run("empty ingredient name", func(t *testing.T) {
		storage := NewMemoryStorage()
		ingredient := ""

		_, err := storage.Create(ingredient)
		if err != ErrIngredientNameCannotBeEmpty {
			t.Errorf("expected %v, got %v", ErrIngredientNameCannotBeEmpty, err)
		}
	})

	t.Run("whitespace normalization", func(t *testing.T) {
		storage := NewMemoryStorage()
		name := "  ChickeN   breast   "
		normalizedName := "chicken breast"

		ingredient, err := storage.Create(name)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ingredient.Name != normalizedName {
			t.Errorf("expected normalized name %q, got %q", normalizedName, ingredient.Name)
		}
	})
}
