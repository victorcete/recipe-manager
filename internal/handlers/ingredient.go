package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"learn-go/internal/storage"
)

// IngredientHandler handles HTTP requests for ingredients
type IngredientHandler struct {
	storage *storage.MemoryStorage
}

// NewIngredientHandler creates a new ingredient handler
func NewIngredientHandler(storage *storage.MemoryStorage) *IngredientHandler {
	return &IngredientHandler{
		storage: storage,
	}
}

// CreateIngredientRequest represents the JSON payload for creating an ingredient
type CreateIngredientRequest struct {
	Name string `json:"name"`
}

// CreateIngredient handles POST /ingredients
func (h *IngredientHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateIngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the ingredient name
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Ingredient name cannot be empty", http.StatusBadRequest)
		return
	}

	// Create the ingredient
	ingredient, err := h.storage.Create(req.Name)
	if err != nil {
		http.Error(w, "Failed to create ingredient", http.StatusInternalServerError)
		return
	}

	// Return the created ingredient as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ingredient)
}
