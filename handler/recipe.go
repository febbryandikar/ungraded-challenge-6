package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"ungraded-challenge-6/entity"
)

type NewRecipeHandler struct {
	*sql.DB
}

func (h *NewRecipeHandler) CreateNewRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var recipe entity.Recipe

	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Failed to create recipe",
		})
		log.Println("Failed to create recipe")
		return
	}

	result, err := h.Exec(`INSERT INTO recipes (recipe_name, description, cook_time, rating) VALUES (?, ?, ?, ?)`, recipe.Name, recipe.Description, recipe.CookTime, recipe.Rating)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error inserting recipes",
		})
		log.Println("Error inserting recipes:", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error inserting orders",
		})
		log.Println("Error inserting orders:", err)
		return
	}
	recipe.ID = int(id)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "Order created successfully",
		Data:    recipe,
	})
}

func (h *NewRecipeHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var recipes []entity.Recipe

	rows, err := h.Query(`SELECT recipe_id, recipe_name, description, cook_time, rating FROM recipes`)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error fetching recipes",
		})
		log.Println("Error fetching recipes:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var recipe entity.Recipe

		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.CookTime, &recipe.Rating)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Message: "Error scanning recipes",
			})
			log.Println("Error scanning recipes:", err)
			return
		}

		recipes = append(recipes, recipe)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Recipes fetched successfully",
		Data:    recipes,
	})
}

func (h *NewRecipeHandler) GetRecipeById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var recipe entity.Recipe
	paramsId := p.ByName("id")

	row, err := h.Query(`SELECT recipe_id, recipe_name, description, cook_time, rating FROM recipes WHERE recipe_id = ?`, paramsId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error fetching recipes",
		})
		log.Println("Error fetching recipes:", err)
		return
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.CookTime, &recipe.Rating)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Message: "Error scanning recipes",
			})
			log.Println("Error scanning recipes:", err)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Recipes not found",
		})
		log.Println("Recipes not found")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Recipes fetched successfully",
		Data:    recipe,
	})
}

func (h *NewRecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var recipe entity.Recipe
	paramsId := p.ByName("id")

	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error updating recipe",
		})
		log.Println("Error updating recipe:", err)
		return
	}

	id, err := strconv.Atoi(paramsId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error fetching orders",
		})
		log.Println("Error fetching orders:", err)
		return
	}
	recipe.ID = id

	result, err := h.Exec(`UPDATE recipes SET recipe_name = ?, description = ?, cook_time = ?, rating = ? WHERE recipe_id = ?`, recipe.Name, recipe.Description, recipe.CookTime, recipe.Rating, recipe.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error updating recipes",
		})
		log.Println("Error updating recipes:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error updating recipes",
		})
		log.Println("Error updating recipes:", err)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Recipes not found",
		})
		log.Println("Recipes not found")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Recipes updated successfully",
		Data:    recipe,
	})
}

func (h *NewRecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	paramsId := p.ByName("id")

	result, err := h.Exec(`DELETE FROM recipes WHERE recipe_id = ?`, paramsId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error deleting recipes",
		})
		log.Println("Error deleting recipes:", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error deleting recipes",
		})
		log.Println("Error deleting recipes:", err)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Recipes not found",
		})
		log.Println("Recipes not found")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Recipes deleted successfully",
	})
}
