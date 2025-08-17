package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/victorcete/recipe-manager/internal/storage"
)

func main() {
	// TODO: Add delete_ingredient MCP tool
	// TODO: Future - add Recipe model and recipe management tools
	// TODO: Future - add search_recipes_by_ingredient tool

	ingredientStorage := storage.NewMemoryStorage()
	ingredientStorage.SeedTestData()
	mcpServer := server.NewMCPServer("ingredient-server", "0.1.0")

	// Tools
	createIngredientTool := mcp.NewTool("create_ingredient",
		mcp.WithDescription("Add exactly one ingredient to your collection. Call this tool separately for each ingredient you want to add. Do not try to add multiple ingredients in a single call."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the single ingredient to add (e.g., 'tomato', 'salt', 'chicken breast')"),
		),
	)

	deleteIngredientTool := mcp.NewTool("delete_ingredient",
		mcp.WithDescription("Delete exactly one ingredient from your collection. Call this tool separately for each ingredient you want to delete. Do not try to delete multiple ingredients in a single call."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the single ingredient to delete (e.g., 'tomato', 'salt', 'chicken breast')"),
		),
	)

	listIngredientsTool := mcp.NewTool("list_ingredients",
		mcp.WithDescription("List all existing ingredients from my collection."),
	)

	updateIngredientTool := mcp.NewTool("update_ingredient",
		mcp.WithDescription("Update exactly one ingredient from your collection. Call this tool separately for each ingredient you want to update. Do not try to update multiple ingredients in a single call."),
		mcp.WithString("original_name",
			mcp.Required(),
			mcp.Description("Name of the already-existing single ingredient"),
		),
		mcp.WithString("new_name",
			mcp.Required(),
			mcp.Description("New name for the single ingredient"),
		),
	)

	// Tool handlers
	mcpServer.AddTool(createIngredientTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}

		ingredient, err := ingredientStorage.Create(name)
		if err != nil {
			var errorMsg string
			switch err {
			// user-friendly storage errors.
			case storage.ErrIngredientNameCannotBeEmpty,
				storage.ErrIngredientNameContainsInvalidChars,
				storage.ErrIngredientNameIsTooShort,
				storage.ErrIngredientNameIsTooLong,
				storage.ErrIngredientNameExists:
				errorMsg = "❌ Error: " + err.Error()
			// default catch for database or system errors, etc.
			default:
				errorMsg = "❌ Error: Failed to create ingredient"
			}
			return mcp.NewToolResultText(errorMsg), nil
		}

		successMsg := fmt.Sprintf("✅ Added %s to your ingredients", ingredient.Name)
		return mcp.NewToolResultText(successMsg), nil
	})

	mcpServer.AddTool(deleteIngredientTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}

		err = ingredientStorage.Delete(name)
		if err != nil {
			var errorMsg string
			switch err {
			// user-friendly storage errors.
			case storage.ErrIngredientNameCannotBeEmpty,
				storage.ErrIngredientNameContainsInvalidChars,
				storage.ErrIngredientNotFound,
				storage.ErrIngredientNameIsTooShort,
				storage.ErrIngredientNameIsTooLong,
				storage.ErrIngredientNameExists:
				errorMsg = "❌ Error: " + err.Error()
			// default catch for database or system errors, etc.
			default:
				errorMsg = "❌ Error: Failed to delete ingredient"
			}
			return mcp.NewToolResultText(errorMsg), nil
		}

		successMsg := fmt.Sprintf("✅ Deleted %s from your ingredients", name)
		return mcp.NewToolResultText(successMsg), nil
	})

	mcpServer.AddTool(listIngredientsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ingredients, err := ingredientStorage.List()
		if err != nil {
			return mcp.NewToolResultText("❌ Error: Failed to fetch ingredients"), nil
		}

		if len(ingredients) == 0 {
			return mcp.NewToolResultText("No ingredients found"), nil
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("📋 Your ingredients (%d total):\n", len(ingredients)))
		for i, ingredient := range ingredients {
			result.WriteString(fmt.Sprintf("%d. %s\n", i+1, ingredient.Name))
		}
		return mcp.NewToolResultText(result.String()), nil
	})

	mcpServer.AddTool(updateIngredientTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		originalName, err := request.RequireString("original_name")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		newName, err := request.RequireString("new_name")
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}

		ingredient, err := ingredientStorage.Update(originalName, newName)
		if err != nil {
			var errorMsg string
			switch err {
			// user-friendly storage errors.
			case storage.ErrIngredientNameCannotBeEmpty,
				storage.ErrIngredientNameContainsInvalidChars,
				storage.ErrIngredientNotFound,
				storage.ErrIngredientNameIsTooShort,
				storage.ErrIngredientNameIsTooLong,
				storage.ErrIngredientNameExists:
				errorMsg = "❌ Error: " + err.Error()
			// default catch for database or system errors, etc.
			default:
				errorMsg = "❌ Error: Failed to update ingredient"
			}
			return mcp.NewToolResultText(errorMsg), nil
		}

		successMsg := fmt.Sprintf("✅ Updated ingredient %s to %s", originalName, ingredient.Name)
		return mcp.NewToolResultText(successMsg), nil
	})

	// create server and start listening
	log.Println("Starting MCP server for ingredient management...")
	stdioServer := server.NewStdioServer(mcpServer)

	if err := stdioServer.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}
