package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	_ "github.com/libsql/go-libsql" // Import the libsql driver.

	"github.com/ngoldack/travel/internal/post/domain"
	"github.com/ngoldack/travel/internal/post/port"
)

type LibsqlPostStore struct {
	log *slog.Logger
	db  *sql.DB
}

func NewLibsqlPostStore(_ context.Context, parentLogger *slog.Logger, db *sql.DB) (*LibsqlPostStore, error) {
	log := parentLogger.WithGroup("post_store")

	store := &LibsqlPostStore{
		log: log,
		db:  db,
	}

	return store, nil
}

// Close implements port.PostStore.
func (store *LibsqlPostStore) Close(_ context.Context) error {
	return store.db.Close()
}

// CreatePost implements port.PostStore.
func (store *LibsqlPostStore) CreatePost(ctx context.Context, post *domain.Post) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO posts (id, created_at, title, description) VALUES (?, ?, ?, ?)", post.ID, post.CreatedAt, post.Title, post.Description)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}
		return err
	}

	return nil
}

// DeletePost implements port.PostStore.
func (store *LibsqlPostStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", id.String())
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}
		return err
	}

	return nil
}

// GetPost implements port.PostStore.
func (store *LibsqlPostStore) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	row := store.db.QueryRowContext(ctx, "SELECT id, created_at, title, description FROM posts WHERE id = ?", id.String())
	if err := row.Err(); err != nil {
		return nil, err
	}

	var post domain.Post
	err := row.Scan(&post.ID, &post.CreatedAt, &post.Title, &post.Description)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// ListPosts implements port.PostStore.
func (store *LibsqlPostStore) ListPosts(ctx context.Context, offset int, size int, filter port.PostFiler) ([]*domain.Post, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, created_at, title, description FROM posts"
	args := []interface{}{}
	if len(filter) > 0 {
		query += " WHERE"
		for k, v := range filter {
			query += fmt.Sprintf(" %s = ?", k)
			args = append(args, v)
		}
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, size, offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, errors.Join(err, txErr)
		}
		return nil, err
	}

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.CreatedAt, &post.Title, &post.Description)
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return nil, errors.Join(err, txErr)
			}
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

// UpdatePost implements port.PostStore.
func (store *LibsqlPostStore) UpdatePost(ctx context.Context, post *domain.Post) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE posts SET title = ?, description = ? WHERE id = ?", post.Title, post.Description, post.ID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}
		return err
	}

	return nil
}

var _ port.PostStore = (*LibsqlPostStore)(nil)
