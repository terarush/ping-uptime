package analytics

import (
	"fmt"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/analytics/domain/entity"
	"ping-uptime/modules/analytics/domain/repository"
	"ping-uptime/modules/analytics/domain/service"
	"ping-uptime/modules/analytics/handler"
	"gorm.io/gorm"
)

type Module struct {
	db               *gorm.DB
	logger           *logger.Logger
	analyticsService *service.AnalyticsService
	analyticsHandler *handler.AnalyticsHandler
	event            *bus.EventBus
}

func (m *Module) Name() string {
	return "analytics"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing analytics module")

	analyticsRepo := repository.NewAnalyticsRepositoryImpl()
	m.logger.Debug("Analytics repository initialized")

	m.analyticsService = service.NewAnalyticsService(analyticsRepo, event)
	m.logger.Debug("Analytics service initialized")

	m.analyticsHandler = handler.NewAnalyticsHandler(m.logger, event, m.analyticsService)
	m.logger.Debug("Analytics handler initialized")

	m.logger.Info("Analytics module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering analytics routes at %s/analytics", basePath)
	m.analyticsHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Analytics routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Analytics module uses existing check_records table")
	return nil
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
