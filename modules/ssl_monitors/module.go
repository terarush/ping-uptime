package ssl_monitors

import (
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/modules/ssl_monitors/domain/entity"
	"ping-uptime/modules/ssl_monitors/domain/repository"
	"ping-uptime/modules/ssl_monitors/domain/service"
	"ping-uptime/modules/ssl_monitors/handler"
	monitorRepository "ping-uptime/modules/monitors/domain/repository"
	monitorService "ping-uptime/modules/monitors/domain/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db          *gorm.DB
	logger      *logger.Logger
	sslService  *service.SSLService
	sslHandler  *handler.SSLHandler
	event       *bus.EventBus
}

func (m *Module) Name() string {
	return "ssl_monitors"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing SSL monitors module")

	// Initialize repositories
	sslRepo := repository.NewSSLRepositoryImpl()
	m.logger.Debug("SSL repository initialized")

	// Initialize services
	m.sslService = service.NewSSLService(sslRepo)
	m.logger.Debug("SSL service initialized")

	// Dependencies from monitor module
	monitorRepo := monitorRepository.NewMonitorRepositoryImpl()
	monServ := monitorService.NewMonitorService(monitorRepo, m.event)

	// Initialize handlers
	m.sslHandler = handler.NewSSLHandler(m.logger, m.sslService, monServ)
	m.logger.Debug("SSL handler initialized")

	m.logger.Info("SSL monitors module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering SSL monitor routes at %s/ssl-monitors", basePath)
	m.sslHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("SSL monitor routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering SSL monitors module migrations")
	return m.db.AutoMigrate(&entity.SSLCert{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func NewModule() *Module {
	return &Module{}
}
