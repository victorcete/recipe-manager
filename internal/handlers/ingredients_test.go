package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/victorcete/recipe-manager/internal/models"
	"github.com/victorcete/recipe-manager/internal/storage"
)

func TestCreateIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := NewIngredientHandler(storage)

	t.Run("successful creation", func(t *testing.T) {
		req := models.CreateIngredientRequest{Name: "tomato"}
		body, _ := json.Marshal(req)

		r := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var ingredient models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &ingredient)

		if ingredient.Name != "tomato" {
			t.Errorf("expected name 'tomato', got '%s'", ingredient.Name)
		}
		if ingredient.ID == 0 {
			t.Error("expected non-zero ID")
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients", nil)
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		req := models.CreateIngredientRequest{Name: "   "}
		body, _ := json.Marshal(req)

		r := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		// Create first ingredient
		req1 := models.CreateIngredientRequest{Name: "onion"}
		body1, _ := json.Marshal(req1)
		r1 := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewReader(body1))
		w1 := httptest.NewRecorder()
		handler.CreateIngredient(w1, r1)

		// Try to create duplicate
		req2 := models.CreateIngredientRequest{Name: "onion"}
		body2, _ := json.Marshal(req2)
		r2 := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewReader(body2))
		w2 := httptest.NewRecorder()
		handler.CreateIngredient(w2, r2)

		if w2.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w2.Code)
		}
	})
}

func TestGetIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := NewIngredientHandler(storage)

	// Create test ingredient
	req := &models.CreateIngredientRequest{Name: "pepper"}
	ingredient, _ := storage.CreateIngredient(req)

	t.Run("successful get", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients/1", nil)
		r.URL.Path = "/api/ingredients/1"
		w := httptest.NewRecorder()

		handler.GetIngredient(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &result)

		if result.ID != ingredient.ID {
			t.Errorf("expected ID %d, got %d", ingredient.ID, result.ID)
		}
		if result.Name != "pepper" {
			t.Errorf("expected name 'pepper', got '%s'", result.Name)
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/ingredients/1", nil)
		w := httptest.NewRecorder()

		handler.GetIngredient(w, r)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients/invalid", nil)
		r.URL.Path = "/api/ingredients/invalid"
		w := httptest.NewRecorder()

		handler.GetIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients/999", nil)
		r.URL.Path = "/api/ingredients/999"
		w := httptest.NewRecorder()

		handler.GetIngredient(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("empty path redirects to GetAll", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients/", nil)
		r.URL.Path = "/api/ingredients/"
		w := httptest.NewRecorder()

		handler.GetIngredient(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var ingredients []models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &ingredients)

		if len(ingredients) == 0 {
			t.Error("expected at least one ingredient")
		}
	})
}

func TestGetAllIngredients(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := NewIngredientHandler(storage)

	// Create test ingredients
	names := []string{"apple", "banana", "cherry"}
	for _, name := range names {
		req := &models.CreateIngredientRequest{Name: name}
		storage.CreateIngredient(req)
	}

	t.Run("get all without search", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients", nil)
		w := httptest.NewRecorder()

		handler.GetAllIngredients(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var ingredients []models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &ingredients)

		if len(ingredients) != 3 {
			t.Errorf("expected 3 ingredients, got %d", len(ingredients))
		}
	})

	t.Run("search with query", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients?search=a", nil)
		w := httptest.NewRecorder()

		handler.GetAllIngredients(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var ingredients []models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &ingredients)

		// Should match "apple" and "banana"
		if len(ingredients) != 2 {
			t.Errorf("expected 2 ingredients matching 'a', got %d", len(ingredients))
		}

		// Results should be sorted alphabetically
		if ingredients[0].Name != "apple" {
			t.Errorf("expected first result to be 'apple', got '%s'", ingredients[0].Name)
		}
		if ingredients[1].Name != "banana" {
			t.Errorf("expected second result to be 'banana', got '%s'", ingredients[1].Name)
		}
	})

	t.Run("search with no matches", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients?search=xyz", nil)
		w := httptest.NewRecorder()

		handler.GetAllIngredients(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var ingredients []models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &ingredients)

		if len(ingredients) != 0 {
			t.Errorf("expected 0 ingredients for 'xyz', got %d", len(ingredients))
		}
	})
}

func TestUpdateIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := NewIngredientHandler(storage)

	// Create test ingredient
	req := &models.CreateIngredientRequest{Name: "original"}
	ingredient, _ := storage.CreateIngredient(req)

	t.Run("successful update", func(t *testing.T) {
		newName := "updated"
		updateReq := models.UpdateIngredientRequest{Name: &newName}
		body, _ := json.Marshal(updateReq)

		r := httptest.NewRequest(http.MethodPut, "/api/ingredients/1", bytes.NewReader(body))
		r.URL.Path = "/api/ingredients/1"
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result models.Ingredient
		json.Unmarshal(w.Body.Bytes(), &result)

		if result.Name != "updated" {
			t.Errorf("expected name 'updated', got '%s'", result.Name)
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients/1", nil)
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		updateReq := models.UpdateIngredientRequest{Name: &ingredient.Name}
		body, _ := json.Marshal(updateReq)

		r := httptest.NewRequest(http.MethodPut, "/api/ingredients/invalid", bytes.NewReader(body))
		r.URL.Path = "/api/ingredients/invalid"
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPut, "/api/ingredients/1", bytes.NewReader([]byte("invalid")))
		r.URL.Path = "/api/ingredients/1"
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		emptyName := "   "
		updateReq := models.UpdateIngredientRequest{Name: &emptyName}
		body, _ := json.Marshal(updateReq)

		r := httptest.NewRequest(http.MethodPut, "/api/ingredients/1", bytes.NewReader(body))
		r.URL.Path = "/api/ingredients/1"
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("ingredient not found", func(t *testing.T) {
		newName := "test"
		updateReq := models.UpdateIngredientRequest{Name: &newName}
		body, _ := json.Marshal(updateReq)

		r := httptest.NewRequest(http.MethodPut, "/api/ingredients/999", bytes.NewReader(body))
		r.URL.Path = "/api/ingredients/999"
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		// Create another ingredient
		req2 := &models.CreateIngredientRequest{Name: "existing"}
		existing, _ := storage.CreateIngredient(req2)

		// Try to update first ingredient to have same name
		updateReq := models.UpdateIngredientRequest{Name: &existing.Name}
		body, _ := json.Marshal(updateReq)

		r := httptest.NewRequest(http.MethodPut, "/api/ingredients/1", bytes.NewReader(body))
		r.URL.Path = "/api/ingredients/1"
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})
}

func TestDeleteIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := NewIngredientHandler(storage)

	// Create test ingredient
	req := &models.CreateIngredientRequest{Name: "to_delete"}
	_, _ = storage.CreateIngredient(req)

	t.Run("successful deletion", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/ingredients/1", nil)
		r.URL.Path = "/api/ingredients/1"
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		// Verify it's gone
		r2 := httptest.NewRequest(http.MethodGet, "/api/ingredients/1", nil)
		r2.URL.Path = "/api/ingredients/1"
		w2 := httptest.NewRecorder()
		handler.GetIngredient(w2, r2)

		if w2.Code != http.StatusNotFound {
			t.Error("ingredient should be deleted")
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/ingredients/1", nil)
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/ingredients/invalid", nil)
		r.URL.Path = "/api/ingredients/invalid"
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("ingredient not found", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/ingredients/999", nil)
		r.URL.Path = "/api/ingredients/999"
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
