package setting

import (
	"time"

	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/settings/domain/entity"
	"ping-uptime/modules/settings/domain/repository"
	"ping-uptime/modules/settings/domain/service"
	"ping-uptime/modules/settings/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db             *gorm.DB
	logger         *logger.Logger
	settingService *service.SettingService
	settingHandler *handler.SettingHandler
	event          *bus.EventBus
}

func (m *Module) Name() string {
	return "setting"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing setting module")

	// Initialize repositories
	settingRepo := repository.NewSettingRepositoryImpl()
	m.logger.Debug("Setting repository initialized")

	// Initialize services
	m.settingService = service.NewSettingService(settingRepo)
	m.logger.Debug("Setting service initialized")

	// Initialize handlers
	m.settingHandler = handler.NewSettingHandler(m.logger, m.event, m.settingService)
	m.logger.Debug("Setting handler initialized")

	m.logger.Info("Setting module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering setting routes at %s/settings", basePath)
	m.settingHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Setting routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering setting module migrations")
	if err := m.db.AutoMigrate(&entity.Setting{}); err != nil {
		return err
	}

	// Define default settings
	defaultSettings := []entity.Setting{
		{
			Key:         "system_name",
			Value:       "ping-uptime",
			Description: "Friendly name of the status monitoring application.",
		},
		{
			Key:         "admin_email",
			Value:       "",
			Description: "Global administrator email address to send alert backups.",
		},
		{
			Key:         "allow_registration",
			Value:       "false",
			Description: "Allow new guest account registrations.",
		},
		{
			Key:         "smtp_host",
			Value:       "",
			Description: "SMTP Host for sending email notifications.",
		},
		{
			Key:         "smtp_port",
			Value:       "587",
			Description: "SMTP Port (e.g. 587 or 465).",
		},
		{
			Key:         "smtp_username",
			Value:       "",
			Description: "SMTP Username/Authentication Email.",
		},
		{
			Key:         "smtp_password",
			Value:       "",
			Description: "SMTP Password/App Password.",
		},
		{
			Key:         "smtp_sender",
			Value:       "",
			Description: "SMTP Sender Address.",
		},
		{
			Key:         "smtp_encryption",
			Value:       "TLS",
			Description: "SMTP Encryption: SSL, TLS, or None.",
		},
	}

	// Push default settings if they do not exist
	for _, s := range defaultSettings {
		var count int64
		if err := m.db.Model(&entity.Setting{}).Where("key = ?", s.Key).Count(&count).Error; err != nil {
			m.logger.Error("Failed to check setting existence for key %s: %v", s.Key, err)
			continue
		}
		if count == 0 {
			now := time.Now()
			s.CreatedAt = now
			s.UpdatedAt = now
			if err := m.db.Create(&s).Error; err != nil {
				m.logger.Error("Failed to create default setting for key %s: %v", s.Key, err)
			} else {
				m.logger.Info("Default setting created: %s = %s", s.Key, s.Value)
			}
		}
	}

	return nil
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
