package storage

import (
	"sync"

	"learn-go/internal/models"
)

// MemoryStorage provides in-memory storage for ingredients
type MemoryStorage struct {
	mu          sync.RWMutex
	ingredients map[int]*models.Ingredient
	nextID      int
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ingredients: make(map[int]*models.Ingredient),
		nextID:      1,
	}
}

// Create adds a new ingredient and returns it with an assigned ID
func (s *MemoryStorage) Create(name string) (*models.Ingredient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ing := models.NewIngredient(s.nextID, name)
	s.ingredients[s.nextID] = ing
	s.nextID++

	return ing, nil
}
