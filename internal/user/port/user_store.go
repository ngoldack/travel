package port

import (
	"github.com/google/uuid"
	"github.com/ngoldack/travel/internal/user/domain"
)

type UserStore interface {
	GetByEmail(email string) (*domain.User, error)
	GetById(id uuid.UUID) (*domain.User, error)
	Store(user *domain.User) error
	Update(id uuid.UUID, user *domain.User) error
	Delete(id uuid.UUID) error
}
