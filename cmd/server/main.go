package main

import (
	"fmt"
	"log"
	"net/http"

	"learn-go/internal/handlers"
	"learn-go/internal/storage"
)

func main() {
	// Initialize storage
	storage := storage.NewMemoryStorage()

	// Initialize handlers
	ingredientHandler := handlers.NewIngredientHandler(storage)

	// Set up routes
	http.HandleFunc("/ingredients", ingredientHandler.CreateIngredient)

	// Start server
	port := "8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Try: curl -X POST http://localhost:%s/ingredients -H \"Content-Type: application/json\" -d '{\"name\":\"tomato\"}'\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
