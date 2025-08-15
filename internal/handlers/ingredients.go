package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/victorcete/recipe-manager/internal/models"
	"github.com/victorcete/recipe-manager/internal/storage"
)

type IngredientHandler struct {
	storage *storage.MemoryStorage
}

func NewIngredientHandler(storage *storage.MemoryStorage) *IngredientHandler {
	return &IngredientHandler{storage: storage}
}

func (h *IngredientHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateIngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	ingredient, err := h.storage.CreateIngredient(&req)
	if err != nil {
		if err == storage.ErrIngredientNameExists {
			http.Error(w, "Ingredient name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ingredient)
}

func (h *IngredientHandler) GetIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/ingredients/")
	if idStr == "" {
		h.GetAllIngredients(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	ingredient, err := h.storage.GetIngredient(id)
	if err != nil {
		if err == storage.ErrIngredientNotFound {
			http.Error(w, "Ingredient not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredient)
}

func (h *IngredientHandler) GetAllIngredients(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")

	var ingredients []*models.Ingredient
	var err error

	if query != "" {
		ingredients, err = h.storage.SearchIngredients(query)
	} else {
		ingredients, err = h.storage.GetAllIngredients()
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}

func (h *IngredientHandler) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/ingredients/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateIngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	ingredient, err := h.storage.UpdateIngredient(id, &req)
	if err != nil {
		if err == storage.ErrIngredientNotFound {
			http.Error(w, "Ingredient not found", http.StatusNotFound)
			return
		}
		if err == storage.ErrIngredientNameExists {
			http.Error(w, "Ingredient name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredient)
}

func (h *IngredientHandler) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/ingredients/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	err = h.storage.DeleteIngredient(id)
	if err != nil {
		if err == storage.ErrIngredientNotFound {
			http.Error(w, "Ingredient not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
