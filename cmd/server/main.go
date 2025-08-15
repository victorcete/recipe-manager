package main

import (
	"log"
	"net/http"
	"os"

	"github.com/victorcete/recipe-manager/internal/handlers"
	"github.com/victorcete/recipe-manager/internal/storage"
)

func main() {
	storage := storage.NewMemoryStorage()
	ingredientHandler := handlers.NewIngredientHandler(storage)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/ingredients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ingredientHandler.CreateIngredient(w, r)
		case http.MethodGet:
			ingredientHandler.GetAllIngredients(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/ingredients/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientHandler.GetIngredient(w, r)
		case http.MethodPut:
			ingredientHandler.UpdateIngredient(w, r)
		case http.MethodDelete:
			ingredientHandler.DeleteIngredient(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
