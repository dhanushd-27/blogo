package model

type SignUp struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// UpdateUser represents the structure for updating user information
// All fields are optional - use pointers to distinguish between zero values and nil (not provided)
type UpdateUser struct {
	ID       int32   `json:"id" validate:"required"`
	Name     *string `json:"name" validate:"omitempty,min=2,max=32"`
	Email    *string `json:"email" validate:"omitempty,email,min=6,max=32"`
	Password *string `json:"password" validate:"omitempty,min=8,max=32"`
}
