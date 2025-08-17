package storage

import (
	"testing"
	"time"
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

	t.Run("ingredient name is too short", func(t *testing.T) {
		storage := NewMemoryStorage()
		ingredient := "XD"

		_, err := storage.Create(ingredient)
		if err != ErrIngredientNameIsTooShort {
			t.Errorf("expected %v, got %v", ErrIngredientNameIsTooShort, err)
		}
	})

	t.Run("ingredient name is too long", func(t *testing.T) {
		storage := NewMemoryStorage()
		ingredient := "Super-Ultra-Mega-Long-Ingredient-Name-That-Goes-On-Forever"

		_, err := storage.Create(ingredient)
		if err != ErrIngredientNameIsTooLong {
			t.Errorf("expected %v, got %v", ErrIngredientNameIsTooLong, err)
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

	t.Run("character regexp validation", func(t *testing.T) {
		storage := NewMemoryStorage()

		testCases := []struct {
			name        string
			input       string
			shouldError bool
			expectedErr error
		}{
			// valid cases
			{"simple ingredient", "tomato", false, nil},
			{"with apostrophe", "Mom's Sauce", false, nil},
			{"with accent", "jalapeño", false, nil},
			{"with hyphen", "extra-virgin", false, nil},
			{"with number", "7-Spice Blend", false, nil},
			{"mixed case", "Chicken Breast", false, nil},

			// invalid cases
			{"with quotes", `"salt"`, true, ErrIngredientNameContainsInvalidChars},
			{"with at symbol", "tom@to", true, ErrIngredientNameContainsInvalidChars},
			{"with brackets", "salt<script>", true, ErrIngredientNameContainsInvalidChars},
			{"with emoji", "🍅", true, ErrIngredientNameContainsInvalidChars},
			{"with parentheses", "salt (sea)", true, ErrIngredientNameContainsInvalidChars},
			{"with period", "Dr. Pepper", true, ErrIngredientNameContainsInvalidChars},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := storage.Create(tc.input)

				if tc.shouldError {
					if err != tc.expectedErr {
						t.Errorf("expected error %v, got %v", tc.expectedErr, err)
					}
				} else {
					if err != nil {
						t.Errorf("expected no error, got %v", err)
					}
				}
			})
		}
	})
}

func TestUpdateIngredient(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		storage := NewMemoryStorage()

		originalName := "tomato"
		storage.Create(originalName)
		newName := "potato"

		time.Sleep(5 * time.Millisecond)

		ingredient, err := storage.Update(originalName, newName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ingredient.Name != newName {
			t.Errorf("expected name %q, got %q", ingredient.Name, newName)
		}
		if !ingredient.UpdatedAt.After(ingredient.CreatedAt) {
			t.Errorf("UpdatedAt should be after CreatedAt")
		}
	})

	t.Run("ingredient not found", func(t *testing.T) {
		storage := NewMemoryStorage()
		storage.Create("tomato")

		_, err := storage.Update("brotato", "potato")
		if err != ErrIngredientNotFound {
			t.Errorf("expected %v, got %v", ErrIngredientNotFound, err)
		}
	})

	t.Run("validation errors", func(t *testing.T) {
		storage := NewMemoryStorage()
		originalName := "tomato"
		storage.Create(originalName)

		_, err := storage.Update(originalName, "tomato")
		if err != ErrIngredientNameExists {
			t.Errorf("expected %v, got %v", ErrIngredientNameExists, err)
		}

		_, err = storage.Update(originalName, " tomato ")
		if err != ErrIngredientNameExists {
			t.Errorf("expected %v, got %v", ErrIngredientNameExists, err)
		}

		_, err = storage.Update(originalName, "ToMaTo   ")
		if err != ErrIngredientNameExists {
			t.Errorf("expected %v, got %v", ErrIngredientNameExists, err)
		}

		_, err = storage.Update(originalName, "")
		if err != ErrIngredientNameCannotBeEmpty {
			t.Errorf("expected %v, got %v", ErrIngredientNameCannotBeEmpty, err)
		}

		_, err = storage.Update(originalName, "a")
		if err != ErrIngredientNameIsTooShort {
			t.Errorf("expected %v, got %v", ErrIngredientNameIsTooShort, err)
		}

		_, err = storage.Update(originalName, "Super-Ultra-Mega-Long-Ingredient-Name-That-Goes-On-Forever")
		if err != ErrIngredientNameIsTooLong {
			t.Errorf("expected %v, got %v", ErrIngredientNameIsTooLong, err)
		}

		_, err = storage.Update(originalName, "<!!tomato>")
		if err != ErrIngredientNameContainsInvalidChars {
			t.Errorf("expected %v, got %v", ErrIngredientNameContainsInvalidChars, err)
		}
	})

	t.Run("updating a different existing ingredient name", func(t *testing.T) {
		storage := NewMemoryStorage()
		storage.Create("tomato")
		storage.Create("basil")

		_, err := storage.Update("tomato", "basil")
		if err != ErrIngredientNameExists {
			t.Errorf("expected %v, got %v", ErrIngredientNameExists, err)
		}
	})
}

func TestListIngredients(t *testing.T) {
	t.Run("list returns valid count and names", func(t *testing.T) {
		storage := NewMemoryStorage()

		storage.Create("tomato")
		storage.Create("basil")
		storage.Create("cheese")
		expectedItems := 3
		expectedNames := []string{"tomato", "basil", "cheese"}

		results, _ := storage.List()

		if len(results) != expectedItems {
			t.Errorf("expected %d ingredients, got %d", expectedItems, len(results))
		}

		// build a map of resulted ingredient names
		foundNames := make(map[string]bool)
		for _, ingredient := range results {
			foundNames[ingredient.Name] = true
		}

		// check each expected name exists
		for _, expectedName := range expectedNames {
			if !foundNames[expectedName] {
				t.Errorf("expected ingredient %q not found", expectedName)
			}
		}
	})

	t.Run("list returns valid IDs", func(t *testing.T) {
		storage := NewMemoryStorage()

		storage.Create("tomato")
		storage.Create("basil")
		storage.Create("cheese")
		expectedIDs := []int{1, 2, 3}

		results, _ := storage.List()

		// build a map of resulted ingredient IDs
		foundIDs := make(map[int]bool)
		for _, ingredient := range results {
			foundIDs[ingredient.ID] = true
		}

		// check each expected ID exists
		for _, expectedID := range expectedIDs {
			if !foundIDs[expectedID] {
				t.Errorf("expected ingredient %d not found", expectedID)
			}
		}
	})

	t.Run("list returns valid timestamps", func(t *testing.T) {
		storage := NewMemoryStorage()

		storage.Create("tomato")
		storage.Create("basil")
		storage.Create("cheese")

		results, _ := storage.List()

		// check all timestamps are not zero
		for _, ingredient := range results {
			if ingredient.CreatedAt.IsZero() {
				t.Errorf("ingredient %q CreatedAt value is zero", ingredient.Name)
			}
			if ingredient.UpdatedAt.IsZero() {
				t.Errorf("ingredient %q UpdatedAt value is zero", ingredient.Name)
			}
		}

	})

	t.Run("empty storage returns empty slice", func(t *testing.T) {
		storage := NewMemoryStorage()
		results, _ := storage.List()

		if len(results) != 0 {
			t.Errorf("expected empty storage, got %v", results)
		}

		if results == nil {
			t.Errorf("expected non-nil slice, got %v", results)
		}
	})
}
