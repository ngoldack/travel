package s3

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image/jpeg"
	"log/slog"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ngoldack/travel/internal/image/domain"
	"github.com/ngoldack/travel/internal/image/port"
)

type ImageStore struct {
	log    *slog.Logger
	db     *sql.DB
	mc     *minio.Client
	bucket string
}

type ImageStoreOption struct {
	Log      *slog.Logger
	Endpoint string
	ID       string
	Secret   string
	Token    string
	Bucket   string
}

// Get implements port.PostImageStore.
func (s *ImageStore) Get(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var i domain.Image
	var url string
	err = tx.QueryRowContext(ctx, "SELECT id, file_name, created_at, url FROM images WHERE id = ?", id.String()).Scan(&i.ID, &i.FileName, &i.CreatedAt, &url)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, errors.Join(err, txErr)
		}
		return nil, err
	}

	obj, err := s.mc.GetObject(ctx, s.bucket, i.FileName, minio.GetObjectOptions{})
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, errors.Join(err, txErr)
		}
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, errors.Join(err, txErr)
		}
		return nil, err
	}

	i.Image, err = jpeg.Decode(buf)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return nil, errors.Join(err, txErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// Store implements port.PostImageStore.
func (s *ImageStore) Store(ctx context.Context, i *domain.Image) error {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, i.Image, nil)
	if err != nil {
		return err
	}

	_, err = s.mc.PutObject(ctx, s.bucket, i.FileName, bytes.NewReader(buf.Bytes()), int64(buf.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/%s/%s", s.mc.EndpointURL().Host, s.bucket, i.FileName)
	_, err = tx.ExecContext(ctx, "INSERT INTO images (id, file_name, created_at, url) VALUES (?, ?, ?, ?)", i.ID, i.FileName, i.CreatedAt, url)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *ImageStore) Close(_ context.Context) error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewImageStore(ctx context.Context, db *sql.DB, opts ImageStoreOption) (*ImageStore, error) {
	log := opts.Log

	log.DebugContext(ctx, "connecting to s3...")
	mc, err := minio.New(opts.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(opts.ID, opts.Secret, opts.Token),
	})
	if err != nil {
		return nil, err
	}
	log.DebugContext(ctx, "connected to s3!")

	is := &ImageStore{
		log: opts.Log,
		db:  db,
		mc:  mc,
	}

	return is, nil
}

var _ port.ImageStore = (*ImageStore)(nil)
