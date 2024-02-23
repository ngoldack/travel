package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Key    string
	UserID uuid.UUID

	createdAt time.Time
	expiresAt time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

func (s *Session) Extend(dur time.Duration) error {
	if s.IsExpired() {
		return ErrSessionExpired
	}

	s.expiresAt = time.Now().Add(dur)
	return nil
}

type KeyFunc func(ctx context.Context) (string, time.Time, error)

func NewSession(ctx context.Context, keyfunc KeyFunc, userID uuid.UUID, expiresAt time.Time) (*Session, error) {
	key, cr, err := keyfunc(ctx)
	if err != nil {
		return nil, err
	}

	return &Session{
		Key:    key,
		UserID: userID,

		createdAt: cr,
		expiresAt: expiresAt,
	}, nil
}
