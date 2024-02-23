package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/ngoldack/travel/internal/validator"
)

type User struct {
	ID uuid.UUID

	Name  string
	Email string
}

type NewUserArgs struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func NewUser(ctx context.Context, args *NewUserArgs) (*User, error) {
	err := validator.Get().StructCtx(ctx, args)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    uuid.New(),
		Name:  args.Name,
		Email: args.Email,
	}, nil
}
