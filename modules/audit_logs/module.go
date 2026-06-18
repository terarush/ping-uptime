package audit_logs

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/audit_logs/domain/entity"
	"ping-uptime/modules/audit_logs/domain/repository"
	"ping-uptime/modules/audit_logs/domain/service"
	"ping-uptime/modules/audit_logs/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db     *gorm.DB
	logger *logger.Logger
	svc    *service.AuditLogService
	h      *handler.AuditLogHandler
}

func (m *Module) Name() string {
	return "audit_logs"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log

	repo := repository.NewAuditLogRepositoryImpl()
	m.svc = service.NewAuditLogService(repo)
	m.h = handler.NewAuditLogHandler(m.logger, m.svc)

	// Subscribe to entity CRUD events
	m.h.SubscribeToEvents(event, m.svc)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.h.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.AuditLog{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
