package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/victorcete/recipe-manager/internal/models"
	"github.com/victorcete/recipe-manager/internal/storage"
)

// Mock template for testing
const mockBaseTemplate = `<!DOCTYPE html>
<html><head><title>{{.Title}}</title></head>
<body>{{template "content" .}}</body></html>
{{define "content"}}<div>{{range .Ingredients}}<span>{{.Name}}</span>{{end}}</div>{{end}}`

const mockFormTemplate = `{{define "ingredient-form"}}<form>{{if .Ingredient}}{{.Ingredient.Name}}{{end}}</form>{{end}}`

const mockTableTemplate = `{{define "ingredients-table"}}{{range .}}<div>{{.Name}}</div>{{end}}{{end}}`

func newTestWebHandler(storage *storage.MemoryStorage) *WebHandler {
	tmpl := template.Must(template.New("base.html").Parse(mockBaseTemplate))
	template.Must(tmpl.New("ingredient-form").Parse(mockFormTemplate))
	template.Must(tmpl.New("ingredients-table").Parse(mockTableTemplate))

	return &WebHandler{
		storage:   storage,
		templates: tmpl,
	}
}

func TestWebHandler_IngredientsPage(t *testing.T) {
	storage := storage.NewMemoryStorage()

	// Add test data
	req1 := &models.CreateIngredientRequest{Name: "tomato"}
	req2 := &models.CreateIngredientRequest{Name: "onion"}
	storage.CreateIngredient(req1)
	storage.CreateIngredient(req2)

	handler := newTestWebHandler(storage)

	r := httptest.NewRequest(http.MethodGet, "/ingredients", nil)
	w := httptest.NewRecorder()

	handler.IngredientsPage(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "tomato") {
		t.Error("expected response to contain 'tomato'")
	}
	if !strings.Contains(body, "onion") {
		t.Error("expected response to contain 'onion'")
	}
	// Just check that we got some HTML content
	if !strings.Contains(body, "<html>") {
		t.Error("expected response to contain HTML")
	}
}

func TestWebHandler_SearchIngredients(t *testing.T) {
	storage := storage.NewMemoryStorage()

	// Add test data
	ingredients := []string{"tomato", "cherry tomato", "onion", "garlic"}
	for _, name := range ingredients {
		req := &models.CreateIngredientRequest{Name: name}
		storage.CreateIngredient(req)
	}

	handler := newTestWebHandler(storage)

	t.Run("search with query", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ingredients/search?search=tomato", nil)
		w := httptest.NewRecorder()

		handler.SearchIngredients(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		body := w.Body.String()
		// For search results, we just verify we got a response
		// The actual filtering logic is tested in storage layer
		if len(body) == 0 {
			t.Error("expected non-empty response")
		}
	})

	t.Run("search with empty query", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ingredients/search", nil)
		w := httptest.NewRecorder()

		handler.SearchIngredients(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		body := w.Body.String()
		if len(body) == 0 {
			t.Error("expected non-empty response")
		}
	})
}

func TestWebHandler_NewIngredientForm(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := newTestWebHandler(storage)

	r := httptest.NewRequest(http.MethodGet, "/ingredients/new", nil)
	w := httptest.NewRecorder()

	handler.NewIngredientForm(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestWebHandler_EditIngredientForm(t *testing.T) {
	storage := storage.NewMemoryStorage()

	// Create test ingredient
	req := &models.CreateIngredientRequest{Name: "tomato"}
	_, _ = storage.CreateIngredient(req)

	handler := newTestWebHandler(storage)

	t.Run("valid ingredient", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ingredients/1/edit", nil)
		r.URL.Path = "/ingredients/1/edit"
		w := httptest.NewRecorder()

		handler.EditIngredientForm(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		body := w.Body.String()
		if len(body) == 0 {
			t.Error("expected non-empty response")
		}
	})

	t.Run("invalid ingredient ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ingredients/invalid/edit", nil)
		r.URL.Path = "/ingredients/invalid/edit"
		w := httptest.NewRecorder()

		handler.EditIngredientForm(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("non-existent ingredient", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ingredients/999/edit", nil)
		r.URL.Path = "/ingredients/999/edit"
		w := httptest.NewRecorder()

		handler.EditIngredientForm(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestWebHandler_CreateIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()
	handler := newTestWebHandler(storage)

	t.Run("successful creation", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "tomato")

		r := httptest.NewRequest(http.MethodPost, "/ingredients", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check that ingredient was created
		ingredients, _ := storage.GetAllIngredients()
		if len(ingredients) != 1 {
			t.Errorf("expected 1 ingredient, got %d", len(ingredients))
		}
		if ingredients[0].Name != "tomato" {
			t.Errorf("expected ingredient name 'tomato', got '%s'", ingredients[0].Name)
		}

		// Check that response contains table data
		body := w.Body.String()
		if len(body) == 0 {
			t.Error("expected non-empty response")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "   ")

		r := httptest.NewRequest(http.MethodPost, "/ingredients", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		// Create first ingredient
		req := &models.CreateIngredientRequest{Name: "onion"}
		storage.CreateIngredient(req)

		form := url.Values{}
		form.Add("name", "onion")

		r := httptest.NewRequest(http.MethodPost, "/ingredients", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/ingredients", nil)
		w := httptest.NewRecorder()

		handler.CreateIngredient(w, r)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestWebHandler_UpdateIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()

	// Create test ingredient
	req := &models.CreateIngredientRequest{Name: "tomato"}
	ingredient, _ := storage.CreateIngredient(req)

	handler := newTestWebHandler(storage)

	t.Run("successful update", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "cherry tomato")

		r := httptest.NewRequest(http.MethodPut, "/ingredients/1", strings.NewReader(form.Encode()))
		r.URL.Path = "/ingredients/1"
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check that ingredient was updated
		updated, _ := storage.GetIngredient(ingredient.ID)
		if updated.Name != "cherry tomato" {
			t.Errorf("expected ingredient name 'cherry tomato', got '%s'", updated.Name)
		}

		// Check that response contains updated table data
		body := w.Body.String()
		if len(body) == 0 {
			t.Error("expected non-empty response")
		}
	})

	t.Run("invalid ingredient ID", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "test")

		r := httptest.NewRequest(http.MethodPut, "/ingredients/invalid", strings.NewReader(form.Encode()))
		r.URL.Path = "/ingredients/invalid"
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("non-existent ingredient", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "test")

		r := httptest.NewRequest(http.MethodPut, "/ingredients/999", strings.NewReader(form.Encode()))
		r.URL.Path = "/ingredients/999"
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.UpdateIngredient(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestWebHandler_DeleteIngredient(t *testing.T) {
	storage := storage.NewMemoryStorage()

	// Create test ingredient
	req := &models.CreateIngredientRequest{Name: "tomato"}
	ingredient, _ := storage.CreateIngredient(req)

	handler := newTestWebHandler(storage)

	t.Run("successful deletion", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/ingredients/1", nil)
		r.URL.Path = "/ingredients/1"
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check that ingredient was deleted
		_, err := storage.GetIngredient(ingredient.ID)
		if err == nil {
			t.Error("expected ingredient to be deleted")
		}
	})

	t.Run("invalid ingredient ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/ingredients/invalid", nil)
		r.URL.Path = "/ingredients/invalid"
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("non-existent ingredient", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/ingredients/999", nil)
		r.URL.Path = "/ingredients/999"
		w := httptest.NewRecorder()

		handler.DeleteIngredient(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
