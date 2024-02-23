package markdown

import (
	"bytes"
	"context"
	"embed"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"

	"github.com/ngoldack/travel/internal/post/domain"
	"github.com/ngoldack/travel/internal/post/port"
)

type MDPostLoader struct {
	posts    embed.FS
	goldmark *goldmark.Markdown
}

func NewMDPostLoader(posts embed.FS) *MDPostLoader {
	gm := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)

	return &MDPostLoader{
		posts:    posts,
		goldmark: &gm,
	}
}

// Load implements port.PostLoader.
func (m *MDPostLoader) Load(ctx context.Context) ([]domain.Post, error) {
	dir, err := m.posts.ReadDir("content/posts")
	if err != nil {
		return nil, err
	}

	posts := make([]domain.Post, 0)
	for _, postDir := range dir {
		if postDir.IsDir() {
			continue
		}

		raw, err := m.posts.ReadFile("content/posts/" + postDir.Name())
		if err != nil {
			return nil, err
		}

		pCtx := parser.NewContext()
		var buf bytes.Buffer
		if err := goldmark.Convert(raw, &buf, parser.WithContext(pCtx)); err != nil {
			return nil, err
		}

		post, err := domain.NewPost(ctx, domain.NewPostArg{
			Title: postDir.Name(),
			Body:  buf.Bytes(),
		})
		if err != nil {
			return nil, err
		}

		posts = append(posts, *post)
	}

	return posts, nil
}

var _ port.PostLoader = (*MDPostLoader)(nil)
