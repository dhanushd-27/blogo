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
	Name *string `json:"name" validate:"required,min=2,max=32"`
}

type CreateBlog struct {
	Title   string `json:"title" validate:"required,min=8,max=150"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdateBlog struct {
	Title   *string `json:"title" validate:"omitempty,min=8,max=150"`
	Content *string `json:"content" validate:"omitempty,min=10"`
}
