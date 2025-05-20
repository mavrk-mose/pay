package models

type UserUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UserFilter struct {
	Role   *string `json:"role,omitempty"`
	Active *bool   `json:"active,omitempty"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}
