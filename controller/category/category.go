package category

import (
	"net/http"
	"encoding/json"


	"product/database"
	"product/model/category"
	
)

func GetCategory (w http.ResponseWriter, r *http.Request) {
	rows,err := database.DB.Query("SELECT * FROM categories")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var categories []category.Category
	for rows.Next() {
		var c category.Category
		if err := rows.Scan(&c.IdCategory,&c.IdUser,&c.IdAdminCategory,&c.CreatedAt,&c.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

func AddCategory (w http.ResponseWriter, r *http.Request) {
	var ac category.Category
	if err := json.NewDecoder(r.Body).Decode(&ac); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new category admin
	query := `
		INSERT INTO categories (id_category_admin, id_user, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id_category`

	// Execute the SQL statement
	var id int
	err := database.DB.QueryRow(query, ac.IdAdminCategory, ac.IdUser).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category added successfully",
		"id":      id,
	})
}

