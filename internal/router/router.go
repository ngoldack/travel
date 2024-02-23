package router

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/ngoldack/travel/internal/image"
	"github.com/ngoldack/travel/internal/post"
	"github.com/ngoldack/travel/internal/router/handler"
	"github.com/ngoldack/travel/internal/router/handler/hx"
	"github.com/ngoldack/travel/internal/router/middleware"
	"github.com/ngoldack/travel/views"
)

type Router struct {
	Log           *slog.Logger
	PostService   *post.Service
	ImageService  *image.Service
	SessionStore  sessions.Store
	sessionSecret string
}

func (r *Router) Handler() http.Handler {
	e := echo.New()

	e.HTTPErrorHandler = ErrorHandler(r.Log)

	e.Use(echomw.CORS())
	e.Use(echomw.RequestID())
	e.Use(middleware.LogMiddleware(r.Log))
	e.Use(echomw.Recover())

	tr := &views.TemplRenderer{}
	e.Renderer = tr

	// Static Files
	e.Static("/assets", "assets")

	// Register Handlers
	e.GET("/", handler.GetIndex())
	e.GET("/reisen", handler.GetReisen())

	e.GET("/post/new", handler.PostNewHandler())

	// Htmx Handlers
	e.POST("/hx/post/render", hx.EditorPreviewRendererHandler())
	e.PUT("/hx/dark-mode", hx.DarkModeHandler())

	return e
}

func NewRouter(log *slog.Logger, postService *post.Service, imageService *image.Service, sessionStore sessions.Store, sessionSecret string) *Router {
	return &Router{
		Log:           log,
		PostService:   postService,
		ImageService:  imageService,
		SessionStore:  sessionStore,
		sessionSecret: sessionSecret,
	}
}
