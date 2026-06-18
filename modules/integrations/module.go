package integration

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/integrations/domain/entity"
	"ping-uptime/modules/integrations/domain/repository"
	"ping-uptime/modules/integrations/domain/service"
	"ping-uptime/modules/integrations/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db                  *gorm.DB
	log                 *logger.Logger
	integrationService  *service.IntegrationService
	integrationHandler  *handler.IntegrationHandler
}

func (m *Module) Name() string {
	return "integration"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, _ *bus.EventBus) error {
	m.db = db
	m.log = log

	// Initialize repositories
	integrationRepo := repository.NewIntegrationRepositoryImpl()

	// Initialize services
	m.integrationService = service.NewIntegrationService(integrationRepo)

	// Initialize handlers
	m.integrationHandler = handler.NewIntegrationHandler(m.integrationService)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.integrationHandler.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Integration{})
}

func (m *Module) Logger() *logger.Logger {
	return m.log
}

func NewModule() *Module {
	return &Module{}
}
