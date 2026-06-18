package teams

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/teams/domain/entity"
	"ping-uptime/modules/teams/domain/repository"
	"ping-uptime/modules/teams/domain/service"
	"ping-uptime/modules/teams/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db  *gorm.DB
	log *logger.Logger
	svc *service.TeamService
	h   *handler.TeamHandler
}

func (m *Module) Name() string {
	return "teams"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.log = log

	repo := repository.NewTeamRepositoryImpl()
	m.svc = service.NewTeamService(repo)
	m.h = handler.NewTeamHandler(m.svc, db)

	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.h.RegisterRoutes(e, basePath)
}

func (m *Module) Migrations() error {
	return m.db.AutoMigrate(&entity.Team{}, &entity.TeamMember{})
}

func (m *Module) Logger() *logger.Logger {
	return m.log
}

func NewModule() *Module {
	return &Module{}
}
