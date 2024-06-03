package category

type Category struct {
	IdCategory int `json:"id_category"`
	IdUser int `json:"id_user"`
	IdAdminCategory int `json:"id_category_admin"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}