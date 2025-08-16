package storage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"learn-go/internal/models"
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

	if len(normalizedName) < IngredientNameMinLength {
		return nil, ErrIngredientNameIsTooShort
	}

	if len(normalizedName) > IngredientNameMaxLength {
		return nil, ErrIngredientNameIsTooLong
	}

	if !isValidIngredientName(normalizedName) {
		return nil, ErrIngredientNameContainsInvalidChars
	}

	if s.IngredientNameExists(normalizedName) {
		return nil, ErrIngredientNameExists
	}

	ingredient := models.NewIngredient(s.nextID, normalizedName)
	s.ingredients[s.nextID] = ingredient
	s.nextID++

	return ingredient, nil
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
