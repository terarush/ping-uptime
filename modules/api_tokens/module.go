package api_tokens

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/api_tokens/domain/entity"
	"ping-uptime/modules/api_tokens/domain/repository"
	"ping-uptime/modules/api_tokens/domain/service"
	"ping-uptime/modules/api_tokens/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db            *gorm.DB
	logger        *logger.Logger
	tokenHandler  *handler.ApiTokenHandler
	tokenService  *service.ApiTokenService
	tokenRepo     repository.ApiTokenRepository
}

func (m *Module) Name() string {
	return "api_tokens"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log

	m.logger.Info("Initializing API tokens module")

	m.tokenRepo = repository.NewApiTokenRepositoryImpl()
	m.logger.Debug("API token repository initialized")

	m.tokenService = service.NewApiTokenService(m.tokenRepo)
	m.logger.Debug("API token service initialized")

	m.tokenHandler = handler.NewApiTokenHandler(m.logger, m.tokenService)
	m.logger.Debug("API token handler initialized")

	m.logger.Info("API tokens module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering API tokens routes at %s/api-tokens", basePath)
	m.tokenHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("API tokens routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Running API tokens migrations")
	return m.db.AutoMigrate(&entity.ApiToken{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
