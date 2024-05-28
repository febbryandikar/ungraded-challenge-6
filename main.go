package main

import (
	"fmt"
	"log"
	"os"
	"ungraded-challenge-6/config"
	"ungraded-challenge-6/handler"
	"ungraded-challenge-6/middleware"
)

func main() {
	router, server := config.SetupServer()
	authDB := &handler.NewAuthHandler{DB: config.GetDatabase()}
	recipeDB := &handler.NewRecipeHandler{DB: config.GetDatabase()}

	router.POST("/register", authDB.Register)
	router.POST("/login", authDB.Login)

	router.GET("/recipes", middleware.AuthMiddleware(os.Getenv("NON_REQUIRED_ROLES"), recipeDB.GetAllRecipes))
	router.GET("/recipes/:id", middleware.AuthMiddleware(os.Getenv("NON_REQUIRED_ROLES"), recipeDB.GetRecipeById))
	router.POST("/recipes", middleware.AuthMiddleware(os.Getenv("REQUIRED_ROLES"), recipeDB.CreateNewRecipe))
	router.PUT("/recipes/:id", middleware.AuthMiddleware(os.Getenv("NON_REQUIRED_ROLES"), recipeDB.UpdateRecipe))
	router.DELETE("/recipes/:id", middleware.AuthMiddleware(os.Getenv("REQUIRED_ROLES"), recipeDB.DeleteRecipe))

	fmt.Println("Server running on port :8080")
	log.Fatal(server.ListenAndServe())
}
