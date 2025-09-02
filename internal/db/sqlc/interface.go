package sqlc

import "context"

//go:generate mockery --name=Querier --output=../mocks --outpkg=mocks

// Querier defines all methods that our database queries implement
type Querier interface {
	// User methods
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id int32) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)

	// Blog methods
	CreateBlog(ctx context.Context, arg CreateBlogParams) (Blog, error)
	DeleteBlog(ctx context.Context, id int32) (int64, error)
	GetBlogByID(ctx context.Context, id int32) (Blog, error)
	ListBlogs(ctx context.Context, arg ListBlogsParams) ([]Blog, error)
	ListBlogsByUserID(ctx context.Context, arg ListBlogsByUserIDParams) ([]Blog, error)
	UpdateBlog(ctx context.Context, arg UpdateBlogParams) (Blog, error)
}
