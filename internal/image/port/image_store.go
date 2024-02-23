package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/ngoldack/travel/internal/image/domain"
)

type ImageStore interface {
	Store(ctx context.Context, i *domain.Image) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Image, error)
}
