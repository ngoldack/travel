package port

import (
	"context"

	"github.com/ngoldack/travel/internal/post/domain"
)

type PostLoader interface {
	Load(ctx context.Context) ([]domain.Post, error)
}
