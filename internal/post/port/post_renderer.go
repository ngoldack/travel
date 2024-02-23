package port

import "github.com/ngoldack/travel/internal/post/domain"

type PostRenderer interface {
	Render(p *domain.Post) ([]byte, error)
}
