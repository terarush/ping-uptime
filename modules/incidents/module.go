package incident

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/incidents/domain/entity"
	"ping-uptime/modules/incidents/domain/repository"
	"ping-uptime/modules/incidents/domain/service"
	"ping-uptime/modules/incidents/handler"
	monitorRepository "ping-uptime/modules/monitors/domain/repository"
	monitorService "ping-uptime/modules/monitors/domain/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db              *gorm.DB
	logger          *logger.Logger
	incidentService *service.IncidentService
	incidentHandler *handler.IncidentHandler
	event           *bus.EventBus
}

func (m *Module) Name() string {
	return "incident"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing incident module")

	// Initialize repositories
	incidentRepo := repository.NewIncidentRepositoryImpl()
	m.logger.Debug("Incident repository initialized")

	// Initialize services
	m.incidentService = service.NewIncidentService(incidentRepo)
	m.logger.Debug("Incident service initialized")

	// Dependencies from monitor module
	monitorRepo := monitorRepository.NewMonitorRepositoryImpl()
	monServ := monitorService.NewMonitorService(monitorRepo, m.event)

	// Initialize handlers
	m.incidentHandler = handler.NewIncidentHandler(m.logger, m.event, m.incidentService, monServ)
	m.logger.Debug("Incident handler initialized")

	m.logger.Info("Incident module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering incident routes at %s/incidents", basePath)
	m.incidentHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Incident routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering incident module migrations")
	return m.db.AutoMigrate(&entity.Incident{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
