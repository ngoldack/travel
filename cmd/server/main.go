package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ngoldack/travel/db"
	"github.com/ngoldack/travel/gen/config"
	authadapter "github.com/ngoldack/travel/internal/auth/adapter"
	"github.com/ngoldack/travel/internal/image"
	"github.com/ngoldack/travel/internal/image/adapter/s3"
	"github.com/ngoldack/travel/internal/logger"
	"github.com/ngoldack/travel/internal/post"
	"github.com/ngoldack/travel/internal/post/adapter"
	"github.com/ngoldack/travel/internal/router"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		cancel()
		os.Exit(1)
	}

	cancel()
}

var configFile = flag.String("config", "pkl/local/config.pkl", "path to the config file")

func run(ctx context.Context) error {
	log := slog.Default()

	cfg, err := config.LoadFromPath(ctx, *configFile)
	if err != nil {
		log.ErrorContext(ctx, "failed to load config", slog.Any("err", err))
		return err
	}

	log = logger.GetLogger(ctx, logger.Options{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})

	db, err := db.NewLibsqlDB(ctx, log, cfg.Db.Url, cfg.Db.Token)
	if err != nil {
		log.ErrorContext(ctx, "failed to create db", slog.Any("err", err))
		return err
	}

	pStore, err := adapter.NewLibsqlPostStore(ctx, log.WithGroup("post_store"), db.DB())
	if err != nil {
		log.ErrorContext(ctx, "failed to create post store", slog.Any("err", err))
		return err
	}
	defer pStore.Close(ctx)

	iStore, err := s3.NewImageStore(ctx, db.DB(), s3.ImageStoreOption{
		Log:      log.WithGroup("image_store"),
		ID:       cfg.S3.Id,
		Secret:   cfg.S3.Secret,
		Token:    cfg.S3.Token,
		Bucket:   cfg.S3.Bucket,
		Endpoint: cfg.S3.Endpoint,
	})
	if err != nil {
		log.ErrorContext(ctx, "failed to create image store", slog.Any("err", err))
		return err
	}
	defer iStore.Close(ctx)

	pService := post.NewPostService(ctx, pStore)
	iService := image.NewImageService(ctx, log, iStore)

	sessionStore, err := authadapter.NewLibsqlSessionStore(ctx, log.WithGroup("session_store"), db.DB(), cfg.Auth.MaxAge.GoDuration(), []byte(cfg.Auth.Secret))
	if err != nil {
		log.ErrorContext(ctx, "failed to create session store", slog.Any("err", err))
		return err
	}
	r := router.NewRouter(log, pService, iService, sessionStore, cfg.Auth.Secret)

	h := r.Handler()
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      h,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ErrorContext(ctx, "failed to start server", slog.Any("err", err))
		}
	}()

	addr := server.Addr
	if cfg.App.Env == "development" {
		addr = "http://" + addr
	}
	log.InfoContext(ctx, "Server started successfully", slog.Any("addr", addr))
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.ErrorContext(ctx, "failed to shutdown server", slog.Any("err", err))
		return err
	}

	log.InfoContext(ctx, "Server shutdown successfully")
	return nil
}
