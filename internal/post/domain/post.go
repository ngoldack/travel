package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ngoldack/travel/internal/validator"
)

type Post struct {
	ID        uuid.UUID
	CreatedAt time.Time

	Title       string
	Description string
	Content     []byte
}

type NewPostArg struct {
	Title       string `validate:"required"`
	Description string `validate:"required"`

	Body []byte
}

func NewPost(ctx context.Context, arg NewPostArg) (*Post, error) {
	err := validator.Get().StructCtx(ctx, arg)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Post{
		ID:          id,
		CreatedAt:   time.Unix(id.Time().UnixTime()),
		Title:       arg.Title,
		Description: arg.Description,
	}, nil
}
