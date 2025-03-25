package models

type Role struct {
	ID          int    `json:"id"`
	RoleName    string `json:"username"`
	Permissions string `json:"email"`
}
