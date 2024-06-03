package user

import (
	// "encoding/json"
	"encoding/json"
	"net/http"
	// "strconv"

	"product/database"
	// "github.com/gorilla/mux"
	"product/model/user"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows,err := database.DB.Query(`SELECT id,name,npm,password,role,created_at FROM users`)
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.NPM, &u.Password, &u.Role, &u.CreatedAt); err != nil {
			http.Error(w,err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users,u)

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}