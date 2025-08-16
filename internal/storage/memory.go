package storage

import (
	"errors"
	"strings"
	"sync"

	"learn-go/internal/models"
)

var (
	ErrIngredientNameCannotBeEmpty = errors.New("ingredient name cannot be empty")
	ErrIngredientNameExists        = errors.New("ingredient name already exists")
)

// MemoryStorage provides in-memory storage for ingredients.
type MemoryStorage struct {
	mu          sync.RWMutex
	ingredients map[int]*models.Ingredient
	nextID      int
}

// NewMemoryStorage creates a new in-memory storage instance.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ingredients: make(map[int]*models.Ingredient),
		nextID:      1,
	}
}

// Create adds a new ingredient and returns it with an assigned ID.
func (s *MemoryStorage) Create(name string) (*models.Ingredient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedName := normalizeIngredientName(name)

	if normalizedName == "" {
		return nil, ErrIngredientNameCannotBeEmpty
	}

	if s.IngredientNameExists(normalizedName) {
		return nil, ErrIngredientNameExists
	}

	ingredient := models.NewIngredient(s.nextID, normalizedName)
	s.ingredients[s.nextID] = ingredient
	s.nextID++

	return ingredient, nil
}

func (s *MemoryStorage) IngredientNameExists(name string) bool {
	for _, ingredient := range s.ingredients {
		if strings.EqualFold(ingredient.Name, name) {
			return true
		}
	}
	return false
}

func normalizeIngredientName(name string) string {
	// remove leading and trailing spaces
	trimmed := strings.TrimSpace(name)

	// normalizes combinations of words that might contain multiple spaces
	normalized := strings.Join(strings.Fields(trimmed), " ")

	// always return lowercased value
	return strings.ToLower(normalized)
}
