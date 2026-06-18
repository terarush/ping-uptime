package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"ping-uptime/internal/pkg/bus"
	"ping-uptime/internal/pkg/email"
	"ping-uptime/internal/pkg/logger"
	incidentEntity "ping-uptime/modules/incidents/domain/entity"
	monitorEntity "ping-uptime/modules/monitors/domain/entity"
	"ping-uptime/modules/notifications/domain/entity"
	"ping-uptime/modules/notifications/domain/repository"
	"ping-uptime/modules/notifications/domain/service"
	"ping-uptime/modules/notifications/handler"
	settingEntity "ping-uptime/modules/settings/domain/entity"
	userEntity "ping-uptime/modules/users/domain/entity"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	db                  *gorm.DB
	logger              *logger.Logger
	notificationService *service.NotificationService
	notificationHandler *handler.NotificationHandler
	event               *bus.EventBus

	// Notification log repo
	notificationLogRepo repository.NotificationLogRepository

	// Persistent Discord bot session
	discordMu      sync.Mutex
	discordSession *discordgo.Session
	discordToken   string
}

func (m *Module) Name() string {
	return "notification"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing notification module")

	// Initialize repositories
	channelRepo := repository.NewNotificationRepositoryImpl()
	logRepo := repository.NewNotificationLogRepositoryImpl()
	m.logger.Debug("Notification repositories initialized")
	m.notificationLogRepo = logRepo

	// Initialize services
	m.notificationService = service.NewNotificationService(channelRepo)
	m.logger.Debug("Notification service initialized")

	// Initialize handlers
	m.notificationHandler = handler.NewNotificationHandler(m.logger, m.event, m.notificationService)
	m.logger.Debug("Notification handler initialized")

	// Register event bus subscribers
	event.SubscribeFunc("incident.created", m.HandleIncidentCreated)
	event.SubscribeFunc("incident.resolved", m.HandleIncidentResolved)
	event.SubscribeFunc("setting.saved", m.HandleSettingSaved)

	// Start persistent Discord bot session if token configured
	m.startDiscordBot()

	m.logger.Info("Notification module initialized successfully")
	return nil
}

func (m *Module) startDiscordBot() {
	token := m.getDiscordBotToken()
	if token == "" {
		m.logger.Debug("No Discord bot token configured, skipping persistent session")
		return
	}

	m.discordMu.Lock()
	defer m.discordMu.Unlock()

	// If already running with same token, skip
	if m.discordSession != nil && m.discordToken == token {
		return
	}

	// Close existing session if any
	m.closeDiscordSessionUnsafe()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		m.logger.Error("Failed to create Discord session: %v", err)
		return
	}

	// Set initial presence — will be updated in Ready handler after connect
	dg.Identify.Presence = discordgo.GatewayStatusUpdate{
		Status: "online",
		Since:  0,
	}
	// Register ready handler -- set online status after connect
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		m.logger.Info("Discord bot connected as %s#%s", r.User.Username, r.User.Discriminator)
		m.updateDiscordPresence(s)
	})

	// Reconnect handler — keep session alive
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Resumed) {
		m.logger.Info("Discord bot session resumed")
	})

	err = dg.Open()
	if err != nil {
		m.logger.Error("Failed to open Discord websocket connection: %v", err)
		return
	}

	m.discordSession = dg
	m.discordToken = token
	m.logger.Info("Discord bot persistent session established")
}

func (m *Module) updateDiscordPresence(s *discordgo.Session) {
	count := int64(0)
	m.db.Model(&monitorEntity.Monitor{}).Count(&count)
	err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online",
		AFK:    false,
		Activities: []*discordgo.Activity{
			{
				Name:    "ping-uptime",
				Type:    discordgo.ActivityTypeGame,
				State:   fmt.Sprintf("%d monitors", count),
				Details: "Monitoring servers uptime",
				Assets: discordgo.Assets{
					LargeImageID: "mp:external/logo/https://raw.githubusercontent.com/rafia9005/ping-uptime/main/public/favicon.svg",
					LargeText:    "Ping Uptime",
				},
			},
		},
	})
	if err != nil {
		m.logger.Warn("Failed to set Discord presence: %v", err)
	} else {
		m.logger.Debug("Discord presence updated")
	}
}

func (m *Module) closeDiscordSessionUnsafe() {
	if m.discordSession != nil {
		m.logger.Info("Closing Discord bot session")
		m.discordSession.Close()
		m.discordSession = nil
	}
	m.discordToken = ""
}

func (m *Module) getDiscordBotToken() string {
	var settings []settingEntity.Setting
	if err := m.db.Find(&settings).Error; err != nil {
		m.logger.Error("Failed to fetch settings for Discord token: %v", err)
		return ""
	}
	for _, s := range settings {
		if s.Key == "discord_bot_token" {
			return s.Value
		}
	}
	return ""
}

// HandleSettingSaved listens for setting.saved events to restart Discord bot on token change
func (m *Module) HandleSettingSaved(event bus.Event) {
	setting, ok := event.Payload.(*settingEntity.Setting)
	if !ok {
		return
	}
	if setting.Key == "discord_bot_token" {
		m.logger.Info("Discord bot token changed, restarting session")
		m.startDiscordBot()
	}
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering notification routes at %s/notification-channels", basePath)
	m.notificationHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Notification routes registered successfully")

	// Register notification log routes
	logHandler := handler.NewNotificationLogHandler(m.logger, m.notificationLogRepo)
	logHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Notification log routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering notification module migrations")
	return m.db.AutoMigrate(&entity.NotificationChannel{}, &entity.NotificationLog{})
}

func (m *Module) Logger() *logger.Logger {
	return m.logger
}

func (m *Module) HandleIncidentCreated(event bus.Event) {
	inc, ok := event.Payload.(*incidentEntity.Incident)
	if !ok {
		m.logger.Warn("Failed to cast incident payload")
		return
	}

	var mon monitorEntity.Monitor
	if err := m.db.First(&mon, inc.MonitorID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor for incident alert: %v", err)
		return
	}

	subject := fmt.Sprintf("[ALERT] Monitor %s is DOWN", mon.Name)
	body := fmt.Sprintf(`
			<h3>Monitor Alert: DOWN</h3>
			<p><strong>Monitor Name:</strong> %s</p>
			<p><strong>Target URL:</strong> <a href="%s">%s</a></p>
			<p><strong>Status:</strong> <span style="color: red; font-weight: bold;">DOWN</span></p>
			<p><strong>Error Details:</strong> %s</p>
			<p><strong>Triggered At:</strong> %s</p>
			<p><strong>Response Latency:</strong> %dms</p>
		`, mon.Name, mon.URL, mon.URL, inc.ErrorMessage, inc.CreatedAt.Format("2006-01-02 15:04:05 MST"), inc.Latency)

	// Send email once per incident
	go m.sendEmail(inc, subject, body, "down")

	// Send channel alerts (discord, telegram, webhook, etc.)
	go m.sendChannelAlerts(inc, &mon, true)
}

func (m *Module) HandleIncidentResolved(event bus.Event) {
	inc, ok := event.Payload.(*incidentEntity.Incident)
	if !ok {
		m.logger.Warn("Failed to cast incident payload")
		return
	}

	var mon monitorEntity.Monitor
	if err := m.db.First(&mon, inc.MonitorID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor for incident resolution alert: %v", err)
		return
	}

	resolvedAtStr := "N/A"
	if inc.ResolvedAt != nil {
		resolvedAtStr = inc.ResolvedAt.Format("2006-01-02 15:04:05 MST")
	}

	subject := fmt.Sprintf("[RESOLVED] Monitor %s is UP", mon.Name)
	body := fmt.Sprintf(`
			<h3>Monitor Alert: RESOLVED</h3>
			<p><strong>Monitor Name:</strong> %s</p>
			<p><strong>Target URL:</strong> <a href="%s">%s</a></p>
			<p><strong>Status:</strong> <span style="color: green; font-weight: bold;">UP (RESOLVED)</span></p>
			<p><strong>Resolved At:</strong> %s</p>
			<p><strong>Response Latency:</strong> %dms</p>
		`, mon.Name, mon.URL, mon.URL, resolvedAtStr, inc.Latency)

	// Send email once per resolution
	go m.sendEmail(inc, subject, body, "up")

	// Send channel alerts
	go m.sendChannelAlerts(inc, &mon, false)
}

// sendEmail sends ONE email per incident — to owner only, no admin backup, no spam.
func (m *Module) sendEmail(incident *incidentEntity.Incident, subject, body, status string) {
	var user userEntity.User
	if err := m.db.First(&user, incident.UserID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor owner for email: %v", err)
		return
	}
	if user.Email == "" {
		return
	}

	var settings []settingEntity.Setting
	if err := m.db.Find(&settings).Error; err != nil {
		m.logger.Error("Failed to fetch SMTP settings: %v", err)
		return
	}

	var smtpHost, smtpUsername, smtpPassword, smtpSender, smtpEncryption string
	var smtpPortVal int
	for _, s := range settings {
		switch s.Key {
		case "smtp_host":
			smtpHost = s.Value
		case "smtp_port":
			smtpPortVal, _ = strconv.Atoi(s.Value)
		case "smtp_username":
			smtpUsername = s.Value
		case "smtp_password":
			smtpPassword = s.Value
		case "smtp_sender":
			smtpSender = s.Value
		case "smtp_encryption":
			smtpEncryption = s.Value
		}
	}

	if smtpHost == "" {
		m.logger.Debug("SMTP not configured, skipping email")
		return
	}

	cfg := email.SMTPConfig{
		Host:       smtpHost,
		Port:       smtpPortVal,
		Username:   smtpUsername,
		Password:   smtpPassword,
		Sender:     smtpSender,
		Encryption: smtpEncryption,
	}

	m.logger.Info("Sending %s alert email to %s", status, user.Email)
	err := email.SendEmail(cfg, user.Email, subject, body)
	if err != nil {
		m.logger.Error("Failed to send email to %s: %v", user.Email, err)
	} else {
		m.logger.Info("Email sent to %s", user.Email)
	}
}

func (m *Module) sendChannelAlerts(incident *incidentEntity.Incident, mon *monitorEntity.Monitor, isDown bool) {
	// Fetch all enabled notification channels for the user
	var channels []entity.NotificationChannel
	if err := m.db.Where("user_id = ? AND enabled = ?", incident.UserID, true).Find(&channels).Error; err != nil {
		m.logger.Error("Failed to fetch notification channels for user %d: %v", incident.UserID, err)
		return
	}

	if len(channels) == 0 {
		return
	}

	// Fetch system settings to get global bot tokens
	var settings []settingEntity.Setting
	if err := m.db.Find(&settings).Error; err != nil {
		m.logger.Error("Failed to fetch system settings for channel alerts: %v", err)
		return
	}

	var telegramBotToken string
	for _, s := range settings {
		if s.Key == "telegram_bot_token" {
			telegramBotToken = s.Value
		}
	}

	for _, ch := range channels {
		switch ch.Type {
		case "discord_bot":
			go m.sendDiscordBotAlert(ch, mon, incident, isDown)
		case "telegram":
			if telegramBotToken == "" {
				m.logger.Warn("Telegram bot token not configured globally, skipping telegram channel")
				continue
			}
			go m.sendTelegramBotAlert(ch, telegramBotToken, mon, incident, isDown)
		case "discord":
			go m.sendDiscordWebhookAlert(ch, mon, incident, isDown)
		case "slack":
			go m.sendSlackWebhookAlert(ch, mon, incident, isDown)
		case "webhook":
			go m.sendCustomWebhookAlert(ch, mon, incident, isDown)
			// email channel type is skipped — email already sent once in sendEmail()
		}
	}
}

// sendDiscordBotAlert uses the persistent session instead of creating one per alert.
func (m *Module) sendDiscordBotAlert(ch entity.NotificationChannel, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
	var config struct {
		ChannelID string `json:"channel_id"`
	}
	if err := json.Unmarshal([]byte(ch.Config), &config); err != nil {
		m.logger.Error("Failed to parse Discord channel config: %v", err)
		return
	}

	if config.ChannelID == "" {
		m.logger.Warn("Discord channel_id is empty, skipping alert")
		return
	}

	m.discordMu.Lock()
	session := m.discordSession
	m.discordMu.Unlock()

	if session == nil {
		m.logger.Warn("Discord bot not connected, cannot send message")
		return
	}

	var content string
	if isDown {
		content = fmt.Sprintf("🔴 **ALERT: Monitor %s is DOWN!**\n**URL:** %s\n**Error:** %s\n**Time:** %s\n**Latency:** %dms",
			mon.Name, mon.URL, incident.ErrorMessage, incident.CreatedAt.Format("2006-01-02 15:04:05 MST"), incident.Latency)
	} else {
		resolvedAtStr := "N/A"
		if incident.ResolvedAt != nil {
			resolvedAtStr = incident.ResolvedAt.Format("2006-01-02 15:04:05 MST")
		}
		content = fmt.Sprintf("🟢 **RESOLVED: Monitor %s is UP!**\n**URL:** %s\n**Time:** %s\n**Latency:** %dms",
			mon.Name, mon.URL, resolvedAtStr, incident.Latency)
	}

	_, err := session.ChannelMessageSend(config.ChannelID, content)
	if err != nil {
		m.logger.Error("Failed to send Discord bot message to channel %s: %v", config.ChannelID, err)
	} else {
		m.logger.Info("Discord bot alert sent successfully to channel %s", config.ChannelID)
	}
}

func (m *Module) sendTelegramBotAlert(ch entity.NotificationChannel, token string, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
	var config struct {
		ChatID string `json:"chat_id"`
	}
	if err := json.Unmarshal([]byte(ch.Config), &config); err != nil {
		m.logger.Error("Failed to parse Telegram channel config: %v", err)
		return
	}

	if config.ChatID == "" {
		m.logger.Warn("Telegram chat_id is empty, skipping alert")
		return
	}

	ctx := context.Background()
	opts := []bot.Option{}
	b, err := bot.New(token, opts...)
	if err != nil {
		m.logger.Error("Failed to create Telegram bot: %v", err)
		return
	}

	var content string
	if isDown {
		content = fmt.Sprintf("🔴 <b>ALERT: Monitor %s is DOWN!</b>\n<b>URL:</b> %s\n<b>Error:</b> %s\n<b>Time:</b> %s\n<b>Latency:</b> %dms",
			html.EscapeString(mon.Name), html.EscapeString(mon.URL), html.EscapeString(incident.ErrorMessage), html.EscapeString(incident.CreatedAt.Format("2006-01-02 15:04:05 MST")), incident.Latency)
	} else {
		resolvedAtStr := "N/A"
		if incident.ResolvedAt != nil {
			resolvedAtStr = incident.ResolvedAt.Format("2006-01-02 15:04:05 MST")
		}
		content = fmt.Sprintf("🟢 <b>RESOLVED: Monitor %s is UP!</b>\n<b>URL:</b> %s\n<b>Time:</b> %s\n<b>Latency:</b> %dms",
			html.EscapeString(mon.Name), html.EscapeString(mon.URL), html.EscapeString(resolvedAtStr), incident.Latency)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    config.ChatID,
		Text:      content,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		m.logger.Error("Failed to send Telegram bot message to chat %s: %v", config.ChatID, err)
	} else {
		m.logger.Info("Telegram bot alert sent successfully to chat %s", config.ChatID)
	}
}

func (m *Module) sendDiscordWebhookAlert(ch entity.NotificationChannel, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
	var config struct {
		WebhookURL string `json:"webhook_url"`
	}
	if err := json.Unmarshal([]byte(ch.Config), &config); err != nil {
		m.logger.Error("Failed to parse Discord webhook config: %v", err)
		return
	}
	if config.WebhookURL == "" {
		return
	}

	var content string
	if isDown {
		content = fmt.Sprintf("🔴 **ALERT: Monitor %s is DOWN!**\n**URL:** %s\n**Error:** %s\n**Time:** %s\n**Latency:** %dms",
			mon.Name, mon.URL, incident.ErrorMessage, incident.CreatedAt.Format("2006-01-02 15:04:05 MST"), incident.Latency)
	} else {
		resolvedAtStr := "N/A"
		if incident.ResolvedAt != nil {
			resolvedAtStr = incident.ResolvedAt.Format("2006-01-02 15:04:05 MST")
		}
		content = fmt.Sprintf("🟢 **RESOLVED: Monitor %s is UP!**\n**URL:** %s\n**Time:** %s\n**Latency:** %dms",
			mon.Name, mon.URL, resolvedAtStr, incident.Latency)
	}

	payload := map[string]string{
		"content": content,
	}
	m.sendJSONPost(config.WebhookURL, payload)
}

func (m *Module) sendSlackWebhookAlert(ch entity.NotificationChannel, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
	var config struct {
		WebhookURL string `json:"webhook_url"`
	}
	if err := json.Unmarshal([]byte(ch.Config), &config); err != nil {
		m.logger.Error("Failed to parse Slack webhook config: %v", err)
		return
	}
	if config.WebhookURL == "" {
		return
	}

	var content string
	if isDown {
		content = fmt.Sprintf("🔴 *ALERT: Monitor %s is DOWN!*\n*URL:* %s\n*Error:* %s\n*Time:* %s\n*Latency:* %dms",
			mon.Name, mon.URL, incident.ErrorMessage, incident.CreatedAt.Format("2006-01-02 15:04:05 MST"), incident.Latency)
	} else {
		resolvedAtStr := "N/A"
		if incident.ResolvedAt != nil {
			resolvedAtStr = incident.ResolvedAt.Format("2006-01-02 15:04:05 MST")
		}
		content = fmt.Sprintf("🟢 *RESOLVED: Monitor %s is UP!*\n*URL:* %s\n*Time:* %s\n*Latency:* %dms",
			mon.Name, mon.URL, resolvedAtStr, incident.Latency)
	}

	payload := map[string]string{
		"text": content,
	}
	m.sendJSONPost(config.WebhookURL, payload)
}

func (m *Module) sendCustomWebhookAlert(ch entity.NotificationChannel, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
	var config struct {
		WebhookURL string `json:"webhook_url"`
	}
	if err := json.Unmarshal([]byte(ch.Config), &config); err != nil {
		m.logger.Error("Failed to parse custom webhook config: %v", err)
		return
	}
	if config.WebhookURL == "" {
		return
	}

	status := "UP"
	if isDown {
		status = "DOWN"
	}

	payload := map[string]any{
		"event": "monitor.status_changed",
		"monitor": map[string]any{
			"id":   mon.ID,
			"name": mon.Name,
			"url":  mon.URL,
		},
		"incident": map[string]any{
			"id":            incident.ID,
			"status":        status,
			"error_message": incident.ErrorMessage,
			"latency":       incident.Latency,
			"created_at":    incident.CreatedAt,
		},
	}
	m.sendJSONPost(config.WebhookURL, payload)
}

func (m *Module) sendJSONPost(url string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		m.logger.Error("Failed to marshal webhook payload: %v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		m.logger.Error("Failed to send POST request to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		m.logger.Error("POST request to %s returned status code %d", url, resp.StatusCode)
	}
}

func NewModule() *Module {
	return &Module{}
}
