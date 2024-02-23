package domain

import (
	"bytes"
	"context"
	"image/jpeg"
	"time"

	stdimage "image"

	"github.com/google/uuid"
	"github.com/ngoldack/travel/internal/validator"
)

type Image struct {
	ID uuid.UUID

	FileName  string
	Image     stdimage.Image
	CreatedAt time.Time
}

type NewImageArg struct {
	FileName string `validate:"required"`
	Data     []byte `validate:"required"`
}

func NewImage(ctx context.Context, arg NewImageArg) (*Image, error) {
	err := validator.Get().StructCtx(ctx, arg)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	stdi, err := jpeg.Decode(bytes.NewReader(arg.Data))
	if err != nil {
		return nil, err
	}

	i := &Image{
		ID:        id,
		CreatedAt: time.Unix(id.Time().UnixTime()),

		FileName: arg.FileName,
		Image:    stdi,
	}
	return i, nil
}
