package status_page

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	monitorRepository "ping-uptime/modules/monitors/domain/repository"
	monitorService "ping-uptime/modules/monitors/domain/service"
	"ping-uptime/modules/status_pages/domain/entity"
	"ping-uptime/modules/status_pages/domain/repository"
	"ping-uptime/modules/status_pages/domain/service"
	"ping-uptime/modules/status_pages/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db                *gorm.DB
	logger            *logger.Logger
	statusPageService *service.StatusPageService
	statusPageHandler *handler.StatusPageHandler
	event             *bus.EventBus
}

func (m *Module) Name() string {
	return "status_page"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing status page module")

	// Initialize repositories
	pageRepo := repository.NewStatusPageRepositoryImpl()
	m.logger.Debug("Status page repository initialized")

	// Initialize services
	m.statusPageService = service.NewStatusPageService(pageRepo)
	m.logger.Debug("Status page service initialized")

	// Dependencies from monitors module
	monitorRepo := monitorRepository.NewMonitorRepositoryImpl()
	monServ := monitorService.NewMonitorService(monitorRepo, m.event)

	// Initialize handlers
	m.statusPageHandler = handler.NewStatusPageHandler(m.logger, m.event, m.statusPageService, monServ)
	m.logger.Debug("Status page handler initialized")

	m.logger.Info("Status page module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering status page routes at %s/status-pages", basePath)
	m.statusPageHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Status page routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering status page module migrations")
	return m.db.AutoMigrate(&entity.StatusPage{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
