package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"learn-go/internal/storage"
)

func main() {
	// Initialize the same storage your HTTP server uses
	ingredientStorage := storage.NewMemoryStorage()

	// Create MCP server with name and version
	mcpServer := server.NewMCPServer("ingredient-server", "1.0.0")

	// Create the create_ingredient tool using the helper function
	createIngredientTool := mcp.NewTool("create_ingredient",
		mcp.WithDescription("Add ONE ingredient to your collection. Call this tool separately for each ingredient you want to add. Do not try to add multiple ingredients in a single call."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the single ingredient to add (e.g., 'tomato', 'salt', 'chicken breast')"),
		),
	)

	// Add the tool with its handler
	mcpServer.AddTool(createIngredientTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the ingredient name using the helper method
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}

		// Use your existing storage logic!
		ingredient, err := ingredientStorage.Create(name)
		if err != nil {
			// Handle your storage errors appropriately
			var errorMsg string
			switch err {
			case storage.ErrIngredientNameCannotBeEmpty:
				errorMsg = "❌ Error: Ingredient name cannot be empty"
			case storage.ErrIngredientNameExists:
				errorMsg = fmt.Sprintf("❌ Error: Ingredient '%s' already exists", name)
			default:
				errorMsg = "❌ Error: Failed to create ingredient"
			}
			return mcp.NewToolResultText(errorMsg), nil
		}

		// Return success result
		successMsg := fmt.Sprintf("✅ Successfully added ingredient: %s (ID: %d)", ingredient.Name, ingredient.ID)
		return mcp.NewToolResultText(successMsg), nil
	})

	// Create stdio server and start listening
	log.Println("Starting MCP server for ingredient management...")
	stdioServer := server.NewStdioServer(mcpServer)
	
	if err := stdioServer.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}