package tags

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/tags/domain/entity"
	"ping-uptime/modules/tags/domain/repository"
	"ping-uptime/modules/tags/domain/service"
	"ping-uptime/modules/tags/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	svc *service.TagService
	h   *handler.TagHandler
}

func (m *Module) Name() string {
	return "tags"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.log = log

	repo := repository.NewTagRepositoryImpl()
	m.svc = service.NewTagService(repo)
	m.h = handler.NewTagHandler(m.svc)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.h.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Tag{}, &entity.MonitorTag{})
}

func (m *Module) Logger() *logger.Logger {
	return m.log
}

func NewModule() *Module {
	return &Module{}
}
