package storage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/victorcete/recipe-manager/internal/models"
)

const (
	IngredientNameMaxLength = 48
	IngredientNameMinLength = 3
)

var (
	ingredientNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\s'\-À-ÿ]+$`)

	ErrIngredientNameCannotBeEmpty        = errors.New("ingredient name cannot be empty")
	ErrIngredientNameContainsInvalidChars = errors.New("ingredient name contains one or more invalid characters")
	ErrIngredientNameExists               = errors.New("ingredient name already exists")
	ErrIngredientNameIsTooLong            = fmt.Errorf("ingredient name cannot exceed %d characters long", IngredientNameMaxLength)
	ErrIngredientNameIsTooShort           = fmt.Errorf("ingredient name must be at least %d characters long", IngredientNameMinLength)
	ErrIngredientNotFound                 = errors.New("ingredient not found")
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

	normalizedName, err := s.validateIngredientName(name)
	if err != nil {
		return nil, err
	}

	if s.IngredientNameExists(normalizedName) {
		return nil, ErrIngredientNameExists
	}

	ingredient := models.NewIngredient(s.nextID, normalizedName)
	s.ingredients[s.nextID] = ingredient
	s.nextID++

	return ingredient, nil
}

func (s *MemoryStorage) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedName, err := s.validateIngredientName(name)
	if err != nil {
		return err
	}

	var targetIngredient *models.Ingredient
	for _, ingredient := range s.ingredients {
		if strings.EqualFold(ingredient.Name, normalizedName) {
			targetIngredient = ingredient
			break
		}
	}

	if targetIngredient == nil {
		return ErrIngredientNotFound
	}

	delete(s.ingredients, targetIngredient.ID)

	return nil
}

func (s *MemoryStorage) List() ([]*models.Ingredient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]*models.Ingredient, 0, len(s.ingredients))
	for _, ingredient := range s.ingredients {
		results = append(results, ingredient)
	}

	return results, nil
}

func (s *MemoryStorage) Update(name, newName string) (*models.Ingredient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// normalize both inputs early
	normalizedName, err := s.validateIngredientName(name)
	if err != nil {
		return nil, err
	}

	normalizedNewName, err := s.validateIngredientName(newName)
	if err != nil {
		return nil, err
	}

	var targetIngredient *models.Ingredient
	for _, ingredient := range s.ingredients {
		if strings.EqualFold(ingredient.Name, normalizedName) {
			targetIngredient = ingredient
			break
		}
	}

	if targetIngredient == nil {
		return nil, ErrIngredientNotFound
	}

	if s.IngredientNameExists(normalizedNewName) {
		return nil, ErrIngredientNameExists
	}

	targetIngredient.Name = normalizedNewName
	targetIngredient.UpdatedAt = time.Now()

	return targetIngredient, nil
}

func (s *MemoryStorage) SeedTestData() ([]*models.Ingredient, error) {
	testIngredients := []string{
		"sal",
		"pimienta negra",
		"ajo en polvo",
		"cebolla en polvo",
		"pimentón",
		"comino",
		"orégano",
		"albahaca seca",
		"tomillo",
		"romero",
		"aceite de oliva",
		"aceite vegetal",
		"vinagre blanco",
		"vinagre de manzana",
		"vinagre balsámico",
		"leche",
		"mantequilla",
		"queso parmesano",
		"huevos",
		"yogur natural",
		"pollo",
		"ternera",
		"pescado blanco",
		"atún en lata",
		"judías",
		"cebolla",
		"ajo fresco",
		"tomate",
		"zanahoria",
		"apio",
		"pimiento",
		"patata",
		"limón",
		"arroz",
		"pasta",
		"pan",
		"harina",
		"avena",
		"azúcar",
		"miel",
		"salsa de soja",
		"caldo de pollo",
		"tomate triturado",
		"mostaza",
		"mahonesa",
		"perejil",
		"cilantro",
		"albahaca fresca",
		"levadura",
		"bicarbonato sódico",
	}
	results := make([]*models.Ingredient, 0, len(testIngredients))

	for _, ingredientName := range testIngredients {
		ingredient, err := s.Create(ingredientName)
		if err != nil {
			continue
		}
		results = append(results, ingredient)
	}
	return results, nil
}

func (s *MemoryStorage) IngredientNameExists(name string) bool {
	for _, ingredient := range s.ingredients {
		if strings.EqualFold(ingredient.Name, name) {
			return true
		}
	}
	return false
}

func (s *MemoryStorage) validateIngredientName(name string) (string, error) {
	normalizedName := normalizeIngredientName(name)

	if normalizedName == "" {
		return "", ErrIngredientNameCannotBeEmpty
	}

	if len(normalizedName) < IngredientNameMinLength {
		return "", ErrIngredientNameIsTooShort
	}

	if len(normalizedName) > IngredientNameMaxLength {
		return "", ErrIngredientNameIsTooLong
	}

	if !isValidIngredientName(normalizedName) {
		return "", ErrIngredientNameContainsInvalidChars
	}

	return normalizedName, nil
}

func normalizeIngredientName(name string) string {
	// remove leading and trailing spaces
	trimmed := strings.TrimSpace(name)

	// normalizes combinations of words that might contain multiple spaces
	normalized := strings.Join(strings.Fields(trimmed), " ")

	// always return lowercased value
	return strings.ToLower(normalized)
}

func isValidIngredientName(name string) bool {
	return ingredientNameRegex.MatchString(name)
}
