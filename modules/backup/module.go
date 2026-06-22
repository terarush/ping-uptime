package backup

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/backup/domain/entity"
	"ping-uptime/modules/backup/domain/repository"
	"ping-uptime/modules/backup/domain/service"
	"ping-uptime/modules/backup/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	svc *service.BackupService
	h   *handler.BackupHandler
}

func (m *Module) Name() string {
	return "backup"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.log = log

	repo := repository.NewBackupRepositoryImpl()
	m.svc = service.NewBackupService(repo)
	m.h = handler.NewBackupHandler(m.svc)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.h.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.BackupRecord{})
}

func (m *Module) Logger() *logger.Logger {
	return m.log
}

func NewModule() *Module {
	return &Module{}
}
