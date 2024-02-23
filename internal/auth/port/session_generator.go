package port

import (
	"context"
	"time"
)

type SessionGenerator interface {
	GenerateKey(ctx context.Context) (string, time.Time, error)
	ValidateKey(ctx context.Context, key string) (bool, error)
}
