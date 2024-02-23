package markdown

import (
	"github.com/ngoldack/travel/internal/post/domain"
	"github.com/ngoldack/travel/internal/post/port"
)

type PostRenderer struct{}

// Render implements port.PostRenderer.
func (*PostRenderer) Render(_ *domain.Post) ([]byte, error) {
	panic("unimplemented")
}

var _ port.PostRenderer = (*PostRenderer)(nil)
