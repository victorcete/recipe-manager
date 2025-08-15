package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/victorcete/recipe-manager/internal/models"
	"github.com/victorcete/recipe-manager/internal/storage"
)

type WebHandler struct {
	storage   *storage.MemoryStorage
	templates *template.Template
}

func NewWebHandler(storage *storage.MemoryStorage) *WebHandler {
	templates := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/ingredients/list.html",
		"web/templates/ingredients/form.html",
	))
	return &WebHandler{
		storage:   storage,
		templates: templates,
	}
}

type PageData struct {
	Title       string
	Ingredients []*models.Ingredient
	Ingredient  *models.Ingredient
}

func (h *WebHandler) IngredientsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ingredients, err := h.storage.GetAllIngredients()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:       "Ingredients",
		Ingredients: ingredients,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) SearchIngredients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	w.Header().Set("Content-Type", "text/html")
	h.templates.ExecuteTemplate(w, "ingredients-table", ingredients)
}

func (h *WebHandler) NewIngredientForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := PageData{
		Title: "Add Ingredient",
	}

	w.Header().Set("Content-Type", "text/html")
	h.templates.ExecuteTemplate(w, "ingredient-form", data)
}

func (h *WebHandler) EditIngredientForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/ingredients/")
	idStr = strings.TrimSuffix(idStr, "/edit")

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

	data := PageData{
		Title:      "Edit Ingredient",
		Ingredient: ingredient,
	}

	w.Header().Set("Content-Type", "text/html")
	h.templates.ExecuteTemplate(w, "ingredient-form", data)
}

func (h *WebHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	req := &models.CreateIngredientRequest{Name: name}
	_, err := h.storage.CreateIngredient(req)
	if err != nil {
		if err == storage.ErrIngredientNameExists {
			http.Error(w, "Ingredient name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return updated ingredients list
	ingredients, err := h.storage.GetAllIngredients()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.templates.ExecuteTemplate(w, "ingredients-table", ingredients)
}

func (h *WebHandler) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/ingredients/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	req := &models.UpdateIngredientRequest{Name: &name}
	_, err = h.storage.UpdateIngredient(id, req)
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

	// Return updated ingredients list
	ingredients, err := h.storage.GetAllIngredients()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.templates.ExecuteTemplate(w, "ingredients-table", ingredients)
}

func (h *WebHandler) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/ingredients/")
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

	// Return empty response (HTMX will remove the element)
	w.WriteHeader(http.StatusOK)
}
