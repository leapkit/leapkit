package internal

import (
	"cmp"
	"embed"
	"os"

	"github.com/leapkit/leapkit/core/assets"
	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server"
	"github.com/leapkit/leapkit/core/session"
	"github.com/leapkit/leapkit/template/internal/home"
	"github.com/leapkit/leapkit/template/public"
)

var (
	//go:embed **/*.html **/*.html *.html
	tmpls embed.FS

	// Assets is the manager for the public assets
	// it allows to watch for changes and reload the assets
	// when changes are made.
	Assets = assets.NewManager(public.Files)

	// DatabaseURL to connect and interact with our database instance.
	DatabaseURL = cmp.Or(os.Getenv("DATABASE_URL"), "leapkit.db")

	// DB is the database connection builder function
	// that will be used by the application based on the driver and
	// connection string.
	DB = db.ConnectionFn(DatabaseURL, db.WithDriver("sqlite3"))
)

// AddRoutes mounts the routes for the application,
// it assumes that the base services have been injected
// in the creation of the server instance.
func AddRoutes(r server.Router) error {
	// LeapKit Middleware
	r.Use(session.Middleware(
		cmp.Or(os.Getenv("SESSION_SECRET"), "d720c059-9664-4980-8169-1158e167ae57"),
		cmp.Or(os.Getenv("SESSION_NAME"), "leapkit_session"),
	))

	r.Use(render.Middleware(
		render.TemplateFS(tmpls, "internal"),

		render.WithDefaultLayout("layout.html"),
		render.WithHelpers(render.AllHelpers),
		render.WithHelpers(map[string]any{
			"assetPath": Assets.PathFor,
		}),
	))

	r.HandleFunc("GET /{$}", home.Index)

	// Mounting the assets manager at the end of the routes
	// so that it can serve the public assets.
	r.HandleFunc(Assets.HandlerPattern(), Assets.HandlerFn)

	return nil
}
