package notification

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/email"
	"ping-uptime/internal/pkg/logger"
	incidentEntity "ping-uptime/modules/incidents/domain/entity"
	monitorEntity "ping-uptime/modules/monitors/domain/entity"
	"ping-uptime/modules/notifications/domain/entity"
	"ping-uptime/modules/notifications/domain/repository"
	"ping-uptime/modules/notifications/domain/service"
	"ping-uptime/modules/notifications/handler"
	settingEntity "ping-uptime/modules/settings/domain/entity"
	userEntity "ping-uptime/modules/users/domain/entity"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db                  *gorm.DB
	logger              *logger.Logger
	notificationService *service.NotificationService
	notificationHandler *handler.NotificationHandler
	event               *bus.EventBus
}

func (m *Module) Name() string {
	return "notification"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing notification module")

	// Initialize repositories
	channelRepo := repository.NewNotificationRepositoryImpl()
	m.logger.Debug("Notification repository initialized")

	// Initialize services
	m.notificationService = service.NewNotificationService(channelRepo)
	m.logger.Debug("Notification service initialized")

	// Initialize handlers
	m.notificationHandler = handler.NewNotificationHandler(m.logger, m.event, m.notificationService)
	m.logger.Debug("Notification handler initialized")

	// Register event bus subscribers
	event.SubscribeFunc("incident.created", m.HandleIncidentCreated)
	event.SubscribeFunc("incident.resolved", m.HandleIncidentResolved)

	m.logger.Info("Notification module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering notification routes at %s/notification-channels", basePath)
	m.notificationHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Notification routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering notification module migrations")
	return m.db.AutoMigrate(&entity.NotificationChannel{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func (m *Module) HandleIncidentCreated(event bus.Event) {
	inc, ok := event.Payload.(*incidentEntity.Incident)
	if !ok {
		m.logger.Warn("Failed to cast incident payload")
		return
	}

	var mon monitorEntity.Monitor
	if err := m.db.First(&mon, inc.MonitorID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor for incident alert: %v", err)
		return
	}

	subject := fmt.Sprintf("[ALERT] Monitor %s is DOWN", mon.Name)
	body := fmt.Sprintf(`
		<h3>Monitor Alert: DOWN</h3>
		<p><strong>Monitor Name:</strong> %s</p>
		<p><strong>Target URL:</strong> <a href="%s">%s</a></p>
		<p><strong>Status:</strong> <span style="color: red; font-weight: bold;">DOWN</span></p>
		<p><strong>Error Details:</strong> %s</p>
		<p><strong>Triggered At:</strong> %s</p>
		<p><strong>Response Latency:</strong> %dms</p>
	`, mon.Name, mon.URL, mon.URL, inc.ErrorMessage, inc.CreatedAt.Format("2006-01-02 15:04:05 MST"), inc.Latency)

	go m.sendEmailAlert(inc, subject, body)
}

func (m *Module) HandleIncidentResolved(event bus.Event) {
	inc, ok := event.Payload.(*incidentEntity.Incident)
	if !ok {
		m.logger.Warn("Failed to cast incident payload")
		return
	}

	var mon monitorEntity.Monitor
	if err := m.db.First(&mon, inc.MonitorID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor for incident resolution alert: %v", err)
		return
	}

	resolvedAtStr := "N/A"
	if inc.ResolvedAt != nil {
		resolvedAtStr = inc.ResolvedAt.Format("2006-01-02 15:04:05 MST")
	}

	subject := fmt.Sprintf("[RESOLVED] Monitor %s is UP", mon.Name)
	body := fmt.Sprintf(`
		<h3>Monitor Alert: RESOLVED</h3>
		<p><strong>Monitor Name:</strong> %s</p>
		<p><strong>Target URL:</strong> <a href="%s">%s</a></p>
		<p><strong>Status:</strong> <span style="color: green; font-weight: bold;">UP (RESOLVED)</span></p>
		<p><strong>Resolved At:</strong> %s</p>
		<p><strong>Response Latency:</strong> %dms</p>
	`, mon.Name, mon.URL, mon.URL, resolvedAtStr, inc.Latency)

	go m.sendEmailAlert(inc, subject, body)
}

func (m *Module) sendEmailAlert(incident *incidentEntity.Incident, subject, body string) {
	var settings []settingEntity.Setting
	if err := m.db.Find(&settings).Error; err != nil {
		m.logger.Error("Failed to fetch settings for email alert: %v", err)
		return
	}

	var smtpHost, smtpUsername, smtpPassword, smtpSender, smtpEncryption, adminEmail string
	var smtpPortVal int

	for _, s := range settings {
		switch s.Key {
		case "smtp_host":
			smtpHost = s.Value
		case "smtp_port":
			smtpPortVal, _ = strconv.Atoi(s.Value)
		case "smtp_username":
			smtpUsername = s.Value
		case "smtp_password":
			smtpPassword = s.Value
		case "smtp_sender":
			smtpSender = s.Value
		case "smtp_encryption":
			smtpEncryption = s.Value
		case "admin_email":
			adminEmail = s.Value
		}
	}

	if smtpHost == "" {
		m.logger.Debug("SMTP not configured, skipping email alert")
		return
	}

	cfg := email.SMTPConfig{
		Host:       smtpHost,
		Port:       smtpPortVal,
		Username:   smtpUsername,
		Password:   smtpPassword,
		Sender:     smtpSender,
		Encryption: smtpEncryption,
	}

	// Fetch monitor owner
	var user userEntity.User
	if err := m.db.First(&user, incident.UserID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor owner user: %v", err)
		return
	}

	// Send to owner user email
	if user.Email != "" {
		m.logger.Info("Sending email alert to owner: %s", user.Email)
		err := email.SendEmail(cfg, user.Email, subject, body)
		if err != nil {
			m.logger.Error("Failed to send email alert to owner %s: %v", user.Email, err)
		}
	}

	// Send to backup admin email if it's set and different
	if adminEmail != "" && adminEmail != user.Email {
		m.logger.Info("Sending email alert to admin backup: %s", adminEmail)
		err := email.SendEmail(cfg, adminEmail, subject, body)
		if err != nil {
			m.logger.Error("Failed to send email alert to admin backup %s: %v", adminEmail, err)
		}
	}
}

func NewModule() *Module {
	return &Module{}
}
