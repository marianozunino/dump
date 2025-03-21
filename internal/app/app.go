package app

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/marianozunino/drop/internal/config"
	"github.com/marianozunino/drop/internal/db"
	"github.com/marianozunino/drop/internal/expiration"
	"github.com/marianozunino/drop/internal/handler"
	middie "github.com/marianozunino/drop/internal/middleware"
)

//go:embed favicon.ico
var faviconFS embed.FS

// App represents the application
type App struct {
	server            *echo.Echo
	expirationManager *expiration.ExpirationManager
	config            *config.Config
	db                *db.DB
}

// New creates a new application instance
func New() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	err = setup(cfg)
	if err != nil {
		return nil, err
	}

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize expiration manager: %v", err)
		return nil, err
	}

	expirationManager, err := expiration.NewExpirationManager(cfg, db)
	if err != nil {
		log.Printf("Warning: Failed to initialize expiration manager: %v", err)
		return nil, err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	app := &App{
		server:            e,
		expirationManager: expirationManager,
		config:            cfg,
		db:                db,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middie.SecurityHeaders())

	registerRoutes(e, app)

	return app, nil
}

// Start starts the application
func (a *App) Start() {
	if a.expirationManager != nil {
		a.expirationManager.Start()
	}

	go func() {
		if err := a.server.Start(":8080"); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	log.Println("Server started on :8080")
}

// Stop stops all application services
func (a *App) Stop() {
	if a.expirationManager != nil {
		a.expirationManager.Stop()
	}
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// setup ensures all necessary directories and files exist
func setup(cfg *config.Config) error {
	if err := os.MkdirAll(cfg.UploadPath, 0o755); err != nil {
		return err
	}

	return nil
}

// registerRoutes registers all HTTP routes
func registerRoutes(e *echo.Echo, app *App) {
	var favicon []byte

	h := handler.NewHandler(app.expirationManager, app.config, app.db)

	// Define routes
	e.GET("/", h.HandleHome)
	e.POST("/", h.HandleUpload)
	e.GET("/favicon.ico", func(c echo.Context) error {
		if favicon == nil {
			data, err := faviconFS.ReadFile("favicon.ico")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Could not read favicon")
			}
			favicon = data
		}
		return c.Blob(http.StatusOK, "image/x-icon", favicon)
	})

	e.GET("/:filename", h.HandleFileAccess)
	e.POST("/:filename", h.HandleFileManagement)
}
