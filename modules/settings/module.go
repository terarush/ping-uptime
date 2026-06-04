package setting

import (
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
	return m.db.AutoMigrate(&entity.Setting{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
