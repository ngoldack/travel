package adapter

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/libsql/go-libsql" // libsql driver

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type sessionRow struct {
	id        string
	data      string
	createdAt time.Time
	updatedAt time.Time
	expiresAt time.Time
}

type LibsqlSessionStore struct {
	log *slog.Logger
	db  *sql.DB

	codecs  []securecookie.Codec
	options *sessions.Options
}

func NewLibsqlSessionStore(_ context.Context, log *slog.Logger, db *sql.DB, maxAge time.Duration, keyPairs ...[]byte) (*LibsqlSessionStore, error) {
	store := &LibsqlSessionStore{
		log: log,
		db:  db,
	}

	store.codecs = securecookie.CodecsFromPairs(keyPairs...)
	store.options = &sessions.Options{
		Path:   "/",
		MaxAge: int(maxAge.Seconds()),
	}

	return store, nil
}

func (store *LibsqlSessionStore) Close(_ context.Context) error {
	return store.db.Close()
}

// Get implements sessions.Store.
func (store *LibsqlSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(store, name)
}

// New implements sessions.Store.
func (store *LibsqlSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	ctx := r.Context()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	session := sessions.NewSession(store, name)
	session.Options = &sessions.Options{
		Path:   store.options.Path,
		MaxAge: store.options.MaxAge,
	}
	session.ID = id.String()
	session.IsNew = true
	session.Values["created_at"] = time.Unix(id.Time().UnixTime())

	if cook, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, cook.Value, &session.ID, store.codecs...)
		if err == nil {
			err = store.load(ctx, session)
			if err == nil {
				session.IsNew = false
			}
		}
	}

	return session, nil
}

// Save implements sessions.Store.
func (store *LibsqlSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	ctx := r.Context()

	var err error
	if s.IsNew {
		err = store.insert(ctx, s)
	} else {
		err = store.save(ctx, s)
	}

	if err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(s.Name(), s.Values, store.codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.Name(),
		Value:    encoded,
		Path:     store.options.Path,
		Expires:  time.Now().Add(time.Duration(store.options.MaxAge) * time.Second),
		HttpOnly: true,
	})

	return nil
}

func (store *LibsqlSessionStore) save(ctx context.Context, s *sessions.Session) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	data, err := securecookie.EncodeMulti(s.Name(), s.Values, store.codecs...)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE sessions SET data = ?, updated_at = ? WHERE id = ?`, data, time.Now(), s.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (store *LibsqlSessionStore) insert(ctx context.Context, s *sessions.Session) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	data, err := securecookie.EncodeMulti(s.Name(), s.Values, store.codecs...)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO sessions (id, data, created_at, updated_at, expires_at) VALUES (?, ?, ?, ?, ?)`, s.ID, data, s.Values["created_at"], time.Now(), time.Now().Add(time.Duration(store.options.MaxAge)*time.Second))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (store *LibsqlSessionStore) load(ctx context.Context, s *sessions.Session) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var sRow sessionRow
	row := tx.QueryRowContext(ctx, "SELECT id, data, created_at, updated_at, expires_at FROM sessions WHERE id = ?", s.ID)
	err = row.Scan(&sRow.id, &sRow.data, &sRow.createdAt, &sRow.updatedAt, &sRow.expiresAt)
	if err != nil {
		return err
	}

	if time.Now().After(sRow.expiresAt) {
		return errors.New("session expired")
	}

	err = securecookie.DecodeMulti(s.Name(), sRow.data, &s.Values, store.codecs...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

var _ sessions.Store = (*LibsqlSessionStore)(nil)
