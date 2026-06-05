package monitor

import (
	"context"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/monitors/domain/entity"
	"ping-uptime/modules/monitors/domain/repository"
	"ping-uptime/modules/monitors/domain/service"
	"ping-uptime/modules/monitors/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db             *gorm.DB
	logger         *logger.Logger
	monitorService *service.MonitorService
	monitorHandler *handler.MonitorHandler
	event          *bus.EventBus
}

func (m *Module) Name() string {
	return "monitor"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing monitor module")

	// Initialize repositories
	monitorRepo := repository.NewMonitorRepositoryImpl()
	m.logger.Debug("Monitor repository initialized")

	// Initialize services
	m.monitorService = service.NewMonitorService(monitorRepo, m.event)
	m.logger.Debug("Monitor service initialized")

	// Initialize handlers
	m.monitorHandler = handler.NewMonitorHandler(m.logger, m.event, m.monitorService)
	m.logger.Debug("Monitor handler initialized")

	// Start background checker scheduler
	go m.monitorService.StartScheduler(context.Background())
	m.logger.Info("Monitor background scheduler started")

	m.logger.Info("Monitor module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering monitor routes at %s/monitors", basePath)
	m.monitorHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Monitor routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering monitor module migrations")
	return m.db.AutoMigrate(&entity.Monitor{}, &entity.CheckRecord{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
