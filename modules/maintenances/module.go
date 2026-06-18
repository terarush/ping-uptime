package maintenances

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/maintenances/domain/entity"
	"ping-uptime/modules/maintenances/domain/repository"
	"ping-uptime/modules/maintenances/domain/service"
	"ping-uptime/modules/maintenances/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db     *gorm.DB
	logger *logger.Logger
	svc    *service.MaintenanceService
	h      *handler.MaintenanceHandler
}

func (m *Module) Name() string {
	return "maintenances"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log

	repo := repository.NewMaintenanceRepositoryImpl()
	m.svc = service.NewMaintenanceService(repo)
	m.h = handler.NewMaintenanceHandler(m.logger, m.svc)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.h.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Maintenance{}, &entity.MaintenanceMonitor{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
