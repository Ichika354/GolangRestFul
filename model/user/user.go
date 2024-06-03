package user

import "time"

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NPM         int       `json:"npm"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}
