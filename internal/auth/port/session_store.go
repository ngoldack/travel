package port

import "github.com/ngoldack/travel/internal/auth/domain"

type SessionStore interface {
	Store(session *domain.Session) error
	Get(key string) (*domain.Session, error)
}
