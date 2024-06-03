package categoryAdmin

import (
    "encoding/json"
    "net/http"
    "strconv"

    "product/database"
    "github.com/gorilla/mux"
    "product/model/categoryAdmin"
)

func GetCategoryAdmins(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT id_category_admin, category, icon, created_at FROM category_admins")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var categoryAdmins []categoryAdmin.CategoryAdmin
    for rows.Next() {
        var ca categoryAdmin.CategoryAdmin
        if err := rows.Scan(&ca.IdCategoryAdmin, &ca.Category, &ca.Icon, &ca.CreatedAt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        categoryAdmins = append(categoryAdmins, ca)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categoryAdmins)
}

func AddCategoryAdmin(w http.ResponseWriter, r *http.Request) {
	var aca categoryAdmin.CategoryAdmin
	if err := json.NewDecoder(r.Body).Decode(&aca); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new category admin
	query := `
		INSERT INTO category_admins (category, icon, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id_category_admin`

	// Execute the SQL statement
	var id int
	err := database.DB.QueryRow(query, aca.Category, aca.Icon).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category admin added successfully",
		"id":      id,
	})
}

func UpdateCategoryAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var uca categoryAdmin.CategoryAdmin
	if err := json.NewDecoder(r.Body).Decode(&uca); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		UPDATE category_admins 
		SET category=$1, icon=$2, updated_at=NOW()
		WHERE id_category_admin=$3`

	result, err := database.DB.Exec(query, uca.Category, uca.Icon, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category admin updated successfully",
	})
}

func DeleteCategoryAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting a category admin
	query := `
		DELETE FROM category_admins
		WHERE id_category_admin=$1`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category admin deleted successfully",
	})
}
