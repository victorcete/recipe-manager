package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/victorcete/recipe-manager/internal/storage"
)

func TestCreateIngredient(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		storage := storage.NewMemoryStorage()
		handler := NewIngredientHandler(storage)

		requestData := CreateIngredientRequest{
			Name: "tomato",
		}

		requestBody, err := json.Marshal(requestData)
		if err != nil {
			t.Fatalf("failed to marshal request: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		handler.CreateIngredient(recorder, req)

		if recorder.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, recorder.Code)
		}
	})

	t.Run("invalid http method", func(t *testing.T) {
		storage := storage.NewMemoryStorage()
		handler := NewIngredientHandler(storage)

		requestData := CreateIngredientRequest{
			Name: "tomato",
		}

		requestBody, err := json.Marshal(requestData)
		if err != nil {
			t.Fatalf("failed to marshal request: %v", err)
		}

		req := httptest.NewRequest(http.MethodPut, "/ingredients", bytes.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		handler.CreateIngredient(recorder, req)

		if recorder.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		storage := storage.NewMemoryStorage()
		handler := NewIngredientHandler(storage)

		malformedJSON := `{"name": "tomato}`

		req := httptest.NewRequest(http.MethodPost, "/ingredients", strings.NewReader(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		handler.CreateIngredient(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("invalid empty name", func(t *testing.T) {
		storage := storage.NewMemoryStorage()
		handler := NewIngredientHandler(storage)

		malformedJSON := `{"name": ""}`

		req := httptest.NewRequest(http.MethodPost, "/ingredients", strings.NewReader(malformedJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		handler.CreateIngredient(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("an ingredient already exists", func(t *testing.T) {
		storage := storage.NewMemoryStorage()
		handler := NewIngredientHandler(storage)

		// create the first ingredient
		req1 := CreateIngredientRequest{Name: "tomato"}
		body1, _ := json.Marshal(req1)
		r1 := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewReader(body1))
		r1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		handler.CreateIngredient(w1, r1)

		// attempt to create the duplicated item
		req2 := CreateIngredientRequest{Name: "tomato"}
		body2, _ := json.Marshal(req2)
		r2 := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewReader(body2))
		r2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		handler.CreateIngredient(w2, r2)

		if w2.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w2.Code)
		}
	})
}
