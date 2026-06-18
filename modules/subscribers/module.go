package subscribers

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/subscribers/domain/entity"
	"ping-uptime/modules/subscribers/domain/repository"
	"ping-uptime/modules/subscribers/domain/service"
	"ping-uptime/modules/subscribers/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db     *gorm.DB
	logger *logger.Logger
	svc    *service.SubscriberService
	h      *handler.SubscriberHandler
}

func (m *Module) Name() string {
	return "subscribers"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log

	repo := repository.NewSubscriberRepositoryImpl()
	m.svc = service.NewSubscriberService(repo)
	m.h = handler.NewSubscriberHandler(m.logger, m.svc)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.h.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Subscriber{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
