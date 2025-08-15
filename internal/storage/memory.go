package storage

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/victorcete/recipe-manager/internal/models"
)

var (
	ErrIngredientNotFound   = errors.New("ingredient not found")
	ErrIngredientNameExists = errors.New("ingredient already exists")
)

type MemoryStorage struct {
	mu          sync.RWMutex
	ingredients map[int]*models.Ingredient
	nextID      int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ingredients: make(map[int]*models.Ingredient),
		nextID:      1,
	}
}

func (s *MemoryStorage) CreateIngredient(req *models.CreateIngredientRequest) (*models.Ingredient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ingredientNameExists(req.Name) {
		return nil, ErrIngredientNameExists
	}

	now := time.Now()
	ingredient := &models.Ingredient{
		ID:        s.nextID,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.ingredients[s.nextID] = ingredient
	s.nextID++

	return ingredient, nil
}

func (s *MemoryStorage) GetIngredient(id int) (*models.Ingredient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ingredient, exists := s.ingredients[id]
	if !exists {
		return nil, ErrIngredientNotFound
	}

	return ingredient, nil
}

func (s *MemoryStorage) GetAllIngredients() ([]*models.Ingredient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ingredients := make([]*models.Ingredient, 0, len(s.ingredients))
	for _, ingredient := range s.ingredients {
		ingredients = append(ingredients, ingredient)
	}

	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].ID < ingredients[j].ID
	})

	return ingredients, nil
}

func (s *MemoryStorage) UpdateIngredient(id int, req *models.UpdateIngredientRequest) (*models.Ingredient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ingredient, exists := s.ingredients[id]
	if !exists {
		return nil, ErrIngredientNotFound
	}

	if req.Name != nil {
		if s.ingredientNameExistsExcluding(*req.Name, id) {
			return nil, ErrIngredientNameExists
		}
		ingredient.Name = *req.Name
	}

	ingredient.UpdatedAt = time.Now()

	return ingredient, nil
}

func (s *MemoryStorage) DeleteIngredient(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.ingredients[id]; !exists {
		return ErrIngredientNotFound
	}

	delete(s.ingredients, id)
	return nil
}

func (s *MemoryStorage) SearchIngredients(query string) ([]*models.Ingredient, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return s.GetAllIngredients()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []*models.Ingredient
	for _, ingredient := range s.ingredients {
		if strings.Contains(strings.ToLower(ingredient.Name), query) {
			matches = append(matches, ingredient)
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Name < matches[j].Name
	})

	return matches, nil
}

func (s *MemoryStorage) ingredientNameExists(name string) bool {
	for _, ingredient := range s.ingredients {
		if strings.EqualFold(ingredient.Name, name) {
			return true
		}
	}
	return false
}

func (s *MemoryStorage) ingredientNameExistsExcluding(name string, excludeID int) bool {
	for _, ingredient := range s.ingredients {
		if ingredient.ID != excludeID && strings.EqualFold(ingredient.Name, name) {
			return true
		}
	}
	return false
}
