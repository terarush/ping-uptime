package app

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/config"
	"ping-uptime/internal/pkg/database"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/server"
	_validator "ping-uptime/internal/pkg/validator"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// App represents the application
type App struct {
	db       *gorm.DB
	server   *server.ServerContext
	modules  []Module
	r        *echo.Echo
	logger   *logger.Logger
	staticFS embed.FS
}

// NewApp creates a new application
func NewApp(cfg *logger.Config, staticFS embed.FS) (*App, error) {
	appLogger, err := logger.NewLogger(*cfg, config.GetString("APP_NAME", "ping-uptime"))
	if err != nil {
		return nil, err
	}
	defer appLogger.Sync()
	return &App{
		modules:  make([]Module, 0),
		logger:   appLogger,
		staticFS: staticFS,
	}, nil
}

func (a *App) SetRouter() *echo.Echo {
	return echo.New()
}

// RegisterModule registers a module with the application
func (a *App) RegisterModule(module Module) {
	a.modules = append(a.modules, module)
	a.logger.Info("Registered module: %s", module.Name())
}

// Initialize initializes the application
func (a *App) Initialize() error {
	a.logger.Info("Initializing application...")

	// Make sure local public directory exists
	if err := os.MkdirAll("public", os.ModePerm); err != nil {
		a.logger.Error("Failed to create public directory: %v", err)
	}

	// Initialize database
	var err *error
	a.db, err = a.SetDatabase().OpenDB()
	if err != nil {
		a.logger.Error("Failed to initialize database: %v", err)
		return *err
	}

	// Set database instance for all modules
	database.DB = a.db

	// event bus initialization
	event := bus.NewEventBus()

	// initialize router
	a.r = a.SetRouter()
	a.r.Use(middleware.Logger())
	a.r.Use(middleware.Recover())
	a.r.Use(middleware.CORS())

	// validate request
	a.r.Validator = _validator.NewCustomValidator()

	// Initialize modules
	for _, module := range a.modules {
		a.logger.Info("Initializing module: %s", module.Name())

		// Create module-specific logger
		moduleLogger := a.logger.WithPrefix(module.Name())
		if err := module.Initialize(a.db, moduleLogger, event); err != nil {
			a.logger.Error("Failed to initialize module %s: %v", module.Name(), err)
			return err
		}

		a.logger.Info("Module initialized: %s", module.Name())
	}

	// Run migrations for all modules
	for _, module := range a.modules {
		err := module.Migrations()
		if err != nil {
			a.logger.Error("Failed to run migrations for module %s: %v", module.Name(), err)
		}
		a.logger.Info("Migrations completed for module: %s", module.Name())
	}

	// Initialize HTTP server
	a.server = a.SetServer()


	for _, module := range a.modules {
		a.logger.Info("Registering routes for module: %s", module.Name())
		module.RegisterRoutes(a.r, "/api")
		a.logger.Info("Routes registered for module: %s", module.Name())
	}

	// API 404 handler — must be registered AFTER all module routes.
	// Any /api/* path that didn't match a real route returns JSON 404
	// instead of falling through to the Vue SPA.
	a.r.Any("/api/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "API endpoint not found",
			"path":    c.Request().URL.Path,
			"method":  c.Request().Method,
		})
	})

	// Serve public static folder for uploaded files and external assets
	a.r.Static("/public", "public")

	// SPA handler: serve actual files if they exist in embedded FS,
	// otherwise fall back to index.html so Vue Router handles the path.
	// e.g: /dashboard, /profile, /about → serve static/index.html
	//      /assets/main.js, /favicon.ico → serve the real file
	staticContent, fsErr := fs.Sub(a.staticFS, "static")
	if fsErr != nil {
		a.logger.Error("Failed to get subdirectory in static embed FS: %v", fsErr)
	}
	httpFS := http.FS(staticContent)

	a.r.GET("/*", func(c echo.Context) error {
		urlPath := c.Param("*")

		// If a path is specified, check if it exists in the embedded FS
		if urlPath != "" {
			file, err := httpFS.Open(urlPath)
			if err == nil {
				defer file.Close()
				info, err := file.Stat()
				if err == nil && !info.IsDir() {
					http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), file)
					return nil
				}
			}
		}

		// Otherwise → return index.html, let Vue Router take over
		indexFile, err := httpFS.Open("index.html")
		if err != nil {
			return c.String(http.StatusNotFound, "index.html not found")
		}
		defer indexFile.Close()
		info, err := indexFile.Stat()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), indexFile)
		return nil
	})

	a.server.Handler = a.r

	a.logger.Info("Application initialization completed")

	for _, v := range a.r.Routes() {
		fmt.Printf("PATH: %v | METHOD: %v\n", v.Path, v.Method)
	}

	return nil
}

// Start starts the application
func (a *App) Start() {
	a.logger.Info("Starting server on %s", a.server.Host)
	a.server.Run()
}

// setup database model
func (a *App) SetDatabase() *database.DBModel {
	return &database.DBModel{
		ServerMode:   config.GetString("SERVER_MODE"),
		Name:         config.GetString("DB_NAME", "ping-uptime"),
		ConnLifeTime: config.GetInt("POOL_CONN_LIFETIME", 60),
	}
}

// Setup Web Server
func (a *App) SetServer() *server.ServerContext {
	return &server.ServerContext{
		Host:         ":" + config.GetString("PORT", "8080"),
		ReadTimeout:  time.Duration(config.GetInt("HTTP_TIMEOUT", 60)),
		WriteTimeout: time.Duration(config.GetInt("HTTP_TIMEOUT", 60)),
	}
}
