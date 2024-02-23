package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/ngoldack/travel/internal/post/domain"
)

type PostFiler map[string]interface{}

type PostStore interface {
	// Close closes the database connection.
	Close(ctx context.Context) error
	// CreatePost creates a new post.
	CreatePost(ctx context.Context, post *domain.Post) error
	// GetPost gets a post by its ID.
	GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	// ListPosts lists all posts.
	ListPosts(ctx context.Context, offset, size int, filter PostFiler) ([]*domain.Post, error)
	// UpdatePost updates a post.
	UpdatePost(ctx context.Context, post *domain.Post) error
	// DeletePost deletes a post by its ID.
	DeletePost(ctx context.Context, id uuid.UUID) error
}
