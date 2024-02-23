package post

import (
	"context"

	"github.com/ngoldack/travel/internal/post/port"
)

type Service struct {
	store port.PostStore
}

func NewPostService(_ context.Context, store port.PostStore) *Service {
	return &Service{store: store}
}

func (s *Service) Close(ctx context.Context) error {
	return s.store.Close(ctx)
}
