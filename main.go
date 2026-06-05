package main

import (
	"embed"
	"log"
	"os"
	"ping-uptime/internal/app"
	"ping-uptime/internal/pkg/config"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/modules/auth"
	incident "ping-uptime/modules/incidents"
	monitor "ping-uptime/modules/monitors"
	notification "ping-uptime/modules/notifications"
	setting "ping-uptime/modules/settings"
	user "ping-uptime/modules/users"
	analytics "ping-uptime/modules/analytics"
)

//go:embed static
var staticFS embed.FS

func main() {

	// Load configuration from .env
	if err := config.Initialize(); err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}

	// initialize logger
	logCfg := logger.DefaultConfig()
	logCfg.Level = config.GetString("SERVER_MODE", "info")

	// Start the application
	application, err := app.NewApp(&logCfg, staticFS)
	if err != nil {
		log.Fatalf("Error creating application : %v", err)
		os.Exit(1)
	}

	// Initialize Auth middleware
	jwtSignatureKey := config.GetJWTService()
	middleware.InitializeAuth(jwtSignatureKey)

	// register modules
	application.RegisterModule(user.NewModule())
	application.RegisterModule(auth.NewModule())
	application.RegisterModule(monitor.NewModule())
	application.RegisterModule(incident.NewModule())
	application.RegisterModule(notification.NewModule())
	application.RegisterModule(setting.NewModule())
	application.RegisterModule(analytics.NewModule())

	// initialize the application
	if err := application.Initialize(); err != nil {
		log.Fatalf("Error initializing application : %v", err)
		os.Exit(1)
	}

	// Start the application
	application.Start()
}
