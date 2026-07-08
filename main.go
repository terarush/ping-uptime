package main

// Ping Uptime API
//
// Website monitoring and uptime tracking API.
//
//	 Schemes: http
//	 BasePath: /api
//	 Version: 1.0.0
//	 Host: localhost:3000
//
//	 SecurityDefinitions:
//	 BearerAuth:
//	   type: apiKey
//	   name: Authorization
//	   in: header
//	   description: "JWT access token. Format: Bearer <token>"
//	 AdminAuth:
//	   type: apiKey
//	   name: X-Admin-Role
//	   in: header
//	   description: "Admin role required"
//
// swagger:meta
import (
	"embed"
	"fmt"
	"log"
	"os"
	"ping-uptime/internal/app"
	"ping-uptime/internal/pkg/config"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	auditLogs "ping-uptime/modules/audit_logs"
	"ping-uptime/modules/auth"
	incident "ping-uptime/modules/incidents"
	maintenances "ping-uptime/modules/maintenances"
	monitor "ping-uptime/modules/monitors"
	notification "ping-uptime/modules/notifications"
	setting "ping-uptime/modules/settings"
	subscribers "ping-uptime/modules/subscribers"
	user "ping-uptime/modules/users"
	analytics "ping-uptime/modules/analytics"
	sslMonitors "ping-uptime/modules/ssl_monitors"
	tags "ping-uptime/modules/tags"
	apiTokens "ping-uptime/modules/api_tokens"
	integration "ping-uptime/modules/integrations"
	"ping-uptime/modules/backup"
	teams "ping-uptime/modules/teams"
	statusPage "ping-uptime/modules/status_pages"
)

//go:embed static
//go:embed docs/swagger.json
var staticFS embed.FS

func main() {

	// Load configuration from .env
	if err := config.Initialize(); err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}

	// Handle --init flag: generate systemd service file
	if config.GetBool("INIT", false) {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get executable path: %v", err)
		}

		svc := fmt.Sprintf(`[Unit]
	Description=Ping Uptime - Website Monitoring & Uptime Tracking
	After=network.target

	[Service]
	Type=simple
	User=%s
	ExecStart=%s
	WorkingDirectory=%s
	Restart=on-failure
	RestartSec=5
	StartLimitIntervalSec=60

	[Install]
	WantedBy=multi-user.target
	`, "nobody", exePath, "/opt/ping-uptime")

		path := "/etc/systemd/system/ping-uptime.service"
		if err := os.WriteFile(path, []byte(svc), 0644); err != nil {
			log.Fatalf("Failed to write systemd service file: %v", err)
		}
		fmt.Printf("Systemd service file created: %s\n", path)
		fmt.Println("Run the following commands to enable and start the service:")
		fmt.Println("  sudo systemctl daemon-reload")
		fmt.Println("  sudo systemctl enable ping-uptime")
		fmt.Println("  sudo systemctl start ping-uptime")
		os.Exit(0)
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
	application.RegisterModule(maintenances.NewModule())
	application.RegisterModule(auditLogs.NewModule())
	application.RegisterModule(subscribers.NewModule())
	application.RegisterModule(sslMonitors.NewModule())
	application.RegisterModule(tags.NewModule())
	application.RegisterModule(apiTokens.NewModule())
	application.RegisterModule(integration.NewModule())
	application.RegisterModule(backup.NewModule())
	application.RegisterModule(teams.NewModule())
	application.RegisterModule(statusPage.NewModule())

	// initialize the application
	if err := application.Initialize(); err != nil {
		log.Fatalf("Error initializing application : %v", err)
		os.Exit(1)
	}

	// Start the application
	application.Start()
}
