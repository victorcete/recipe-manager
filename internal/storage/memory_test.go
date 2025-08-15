package storage

import (
	"testing"
	"time"

	"github.com/victorcete/recipe-manager/internal/models"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()
	if storage == nil {
		t.Fatal("NewMemoryStorage() returned nil")
	}
	if storage.ingredients == nil {
		t.Error("ingredients map is nil")
	}
	if storage.nextID != 1 {
		t.Errorf("expected nextID to be 1, got %d", storage.nextID)
	}
}

func TestCreateIngredient(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("successful creation", func(t *testing.T) {
		req := &models.CreateIngredientRequest{Name: "tomato"}
		ingredient, err := storage.CreateIngredient(req)

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
			t.Error("CreatedAt should not be zero")
		}
		if ingredient.UpdatedAt.IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		req1 := &models.CreateIngredientRequest{Name: "onion"}
		req2 := &models.CreateIngredientRequest{Name: "onion"}

		_, err := storage.CreateIngredient(req1)
		if err != nil {
			t.Fatalf("first creation failed: %v", err)
		}

		_, err = storage.CreateIngredient(req2)
		if err != ErrIngredientNameExists {
			t.Errorf("expected ErrIngredientNameExists, got %v", err)
		}
	})

	t.Run("case insensitive duplicate", func(t *testing.T) {
		req1 := &models.CreateIngredientRequest{Name: "Garlic"}
		req2 := &models.CreateIngredientRequest{Name: "garlic"}

		_, err := storage.CreateIngredient(req1)
		if err != nil {
			t.Fatalf("first creation failed: %v", err)
		}

		_, err = storage.CreateIngredient(req2)
		if err != ErrIngredientNameExists {
			t.Errorf("expected ErrIngredientNameExists for case insensitive duplicate, got %v", err)
		}
	})
}

func TestGetIngredient(t *testing.T) {
	storage := NewMemoryStorage()
	req := &models.CreateIngredientRequest{Name: "pepper"}
	created, _ := storage.CreateIngredient(req)

	t.Run("existing ingredient", func(t *testing.T) {
		ingredient, err := storage.GetIngredient(created.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ingredient.ID != created.ID {
			t.Errorf("expected ID %d, got %d", created.ID, ingredient.ID)
		}
		if ingredient.Name != "pepper" {
			t.Errorf("expected name 'pepper', got '%s'", ingredient.Name)
		}
	})

	t.Run("non-existent ingredient", func(t *testing.T) {
		_, err := storage.GetIngredient(999)
		if err != ErrIngredientNotFound {
			t.Errorf("expected ErrIngredientNotFound, got %v", err)
		}
	})
}

func TestGetAllIngredients(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("empty storage", func(t *testing.T) {
		ingredients, err := storage.GetAllIngredients()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ingredients) != 0 {
			t.Errorf("expected 0 ingredients, got %d", len(ingredients))
		}
	})

	t.Run("multiple ingredients sorted by ID", func(t *testing.T) {
		names := []string{"carrot", "apple", "banana"}
		for _, name := range names {
			req := &models.CreateIngredientRequest{Name: name}
			storage.CreateIngredient(req)
		}

		ingredients, err := storage.GetAllIngredients()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ingredients) != 3 {
			t.Errorf("expected 3 ingredients, got %d", len(ingredients))
		}

		// Should be sorted by ID (creation order)
		expectedNames := []string{"carrot", "apple", "banana"}
		for i, ingredient := range ingredients {
			if ingredient.Name != expectedNames[i] {
				t.Errorf("expected name '%s' at index %d, got '%s'", expectedNames[i], i, ingredient.Name)
			}
		}
	})
}

func TestUpdateIngredient(t *testing.T) {
	storage := NewMemoryStorage()
	req := &models.CreateIngredientRequest{Name: "original"}
	created, _ := storage.CreateIngredient(req)
	time.Sleep(time.Millisecond) // Ensure UpdatedAt changes

	t.Run("successful update", func(t *testing.T) {
		newName := "updated"
		updateReq := &models.UpdateIngredientRequest{Name: &newName}
		
		updated, err := storage.UpdateIngredient(created.ID, updateReq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Name != "updated" {
			t.Errorf("expected name 'updated', got '%s'", updated.Name)
		}
		if !updated.UpdatedAt.After(updated.CreatedAt) {
			t.Error("UpdatedAt should be after CreatedAt")
		}
	})

	t.Run("non-existent ingredient", func(t *testing.T) {
		newName := "test"
		updateReq := &models.UpdateIngredientRequest{Name: &newName}
		
		_, err := storage.UpdateIngredient(999, updateReq)
		if err != ErrIngredientNotFound {
			t.Errorf("expected ErrIngredientNotFound, got %v", err)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		// Create another ingredient
		req2 := &models.CreateIngredientRequest{Name: "existing"}
		existing, _ := storage.CreateIngredient(req2)

		// Try to update first ingredient to have same name as second
		updateReq := &models.UpdateIngredientRequest{Name: &existing.Name}
		_, err := storage.UpdateIngredient(created.ID, updateReq)
		
		if err != ErrIngredientNameExists {
			t.Errorf("expected ErrIngredientNameExists, got %v", err)
		}
	})

	t.Run("update to same name", func(t *testing.T) {
		// Should allow updating to the same name
		sameName := created.Name
		updateReq := &models.UpdateIngredientRequest{Name: &sameName}
		
		_, err := storage.UpdateIngredient(created.ID, updateReq)
		if err != nil {
			t.Errorf("should allow updating to same name, got error: %v", err)
		}
	})
}

func TestDeleteIngredient(t *testing.T) {
	storage := NewMemoryStorage()
	req := &models.CreateIngredientRequest{Name: "to_delete"}
	created, _ := storage.CreateIngredient(req)

	t.Run("successful deletion", func(t *testing.T) {
		err := storage.DeleteIngredient(created.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify it's gone
		_, err = storage.GetIngredient(created.ID)
		if err != ErrIngredientNotFound {
			t.Errorf("expected ingredient to be deleted")
		}
	})

	t.Run("non-existent ingredient", func(t *testing.T) {
		err := storage.DeleteIngredient(999)
		if err != ErrIngredientNotFound {
			t.Errorf("expected ErrIngredientNotFound, got %v", err)
		}
	})
}

func TestSearchIngredients(t *testing.T) {
	storage := NewMemoryStorage()
	
	// Create test data
	ingredients := []string{"tomato", "cherry tomato", "tomato paste", "onion", "garlic"}
	for _, name := range ingredients {
		req := &models.CreateIngredientRequest{Name: name}
		storage.CreateIngredient(req)
	}

	t.Run("empty query returns all", func(t *testing.T) {
		results, err := storage.SearchIngredients("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 5 {
			t.Errorf("expected 5 ingredients, got %d", len(results))
		}
		// Should be sorted by ID (same as GetAllIngredients)
		if results[0].Name != "tomato" {
			t.Errorf("expected first ingredient to be 'tomato', got '%s'", results[0].Name)
		}
	})

	t.Run("search with results", func(t *testing.T) {
		results, err := storage.SearchIngredients("tomato")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 3 {
			t.Errorf("expected 3 results for 'tomato', got %d", len(results))
		}
		
		// Should be sorted alphabetically by name
		expectedOrder := []string{"cherry tomato", "tomato", "tomato paste"}
		for i, result := range results {
			if result.Name != expectedOrder[i] {
				t.Errorf("expected '%s' at index %d, got '%s'", expectedOrder[i], i, result.Name)
			}
		}
	})

	t.Run("case insensitive search", func(t *testing.T) {
		results, err := storage.SearchIngredients("TOMATO")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 3 {
			t.Errorf("expected 3 results for case insensitive search, got %d", len(results))
		}
	})

	t.Run("no results", func(t *testing.T) {
		results, err := storage.SearchIngredients("nonexistent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 0 {
			t.Errorf("expected 0 results, got %d", len(results))
		}
	})

	t.Run("whitespace handling", func(t *testing.T) {
		results, err := storage.SearchIngredients("  tomato  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 3 {
			t.Errorf("expected 3 results with whitespace, got %d", len(results))
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	storage := NewMemoryStorage()
	
	// Test concurrent reads and writes
	done := make(chan bool)
	
	// Concurrent writers
	go func() {
		for i := 0; i < 100; i++ {
			req := &models.CreateIngredientRequest{Name: "ingredient" + string(rune(i))}
			storage.CreateIngredient(req)
		}
		done <- true
	}()
	
	// Concurrent readers
	go func() {
		for i := 0; i < 100; i++ {
			storage.GetAllIngredients()
		}
		done <- true
	}()
	
	// Wait for both goroutines
	<-done
	<-done
	
	// Verify final state
	ingredients, err := storage.GetAllIngredients()
	if err != nil {
		t.Fatalf("unexpected error after concurrent access: %v", err)
	}
	
	// Should have some ingredients (exact count may vary due to duplicates)
	if len(ingredients) == 0 {
		t.Error("expected some ingredients after concurrent writes")
	}
}