package image

import (
	"context"
	"log/slog"

	"github.com/ngoldack/travel/internal/image/port"
)

type Service struct {
	log   *slog.Logger
	store port.ImageStore
}

func NewImageService(_ context.Context, log *slog.Logger, store port.ImageStore) *Service {
	return &Service{
		log:   log.WithGroup("image_service"),
		store: store,
	}
}
