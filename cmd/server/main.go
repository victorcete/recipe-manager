package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/victorcete/recipe-manager/internal/handlers"
	"github.com/victorcete/recipe-manager/internal/storage"
)

func main() {
	storage := storage.NewMemoryStorage()
	ingredientHandler := handlers.NewIngredientHandler(storage)
	webHandler := handlers.NewWebHandler(storage)

	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Web UI routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ingredients", http.StatusFound)
	})
	mux.HandleFunc("/ingredients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			webHandler.IngredientsPage(w, r)
		case http.MethodPost:
			webHandler.CreateIngredient(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/ingredients/search", webHandler.SearchIngredients)
	mux.HandleFunc("/ingredients/new", webHandler.NewIngredientForm)

	// Web CRUD operations
	mux.HandleFunc("/ingredients/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") {
			webHandler.EditIngredientForm(w, r)
		} else {
			switch r.Method {
			case http.MethodPost:
				webHandler.CreateIngredient(w, r)
			case http.MethodPut:
				webHandler.UpdateIngredient(w, r)
			case http.MethodDelete:
				webHandler.DeleteIngredient(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// JSON API routes
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

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	addr := ":" + port
	log.Printf("Server started on http://%s:%s", host, port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
