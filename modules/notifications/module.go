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
	m.logger.Debug("Notification repository initialized")

	// Initialize services
	m.notificationService = service.NewNotificationService(channelRepo)
	m.logger.Debug("Notification service initialized")

	// Initialize handlers
	m.notificationHandler = handler.NewNotificationHandler(m.logger, m.event, m.notificationService)
	m.logger.Debug("Notification handler initialized")

	// Register event bus subscribers
	event.SubscribeFunc("incident.created", m.HandleIncidentCreated)
	event.SubscribeFunc("incident.resolved", m.HandleIncidentResolved)

	m.logger.Info("Notification module initialized successfully")
	return nil
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering notification routes at %s/notification-channels", basePath)
	m.notificationHandler.RegisterRoutes(e, basePath)
	m.logger.Debug("Notification routes registered successfully")
}

func (m *Module) Migrations() error {
	m.logger.Info("Registering notification module migrations")
	return m.db.AutoMigrate(&entity.NotificationChannel{})
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

	go m.sendEmailAlert(inc, subject, body)
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

	go m.sendEmailAlert(inc, subject, body)
	go m.sendChannelAlerts(inc, &mon, false)
}

func (m *Module) sendEmailAlert(incident *incidentEntity.Incident, subject, body string) {
	var settings []settingEntity.Setting
	if err := m.db.Find(&settings).Error; err != nil {
		m.logger.Error("Failed to fetch settings for email alert: %v", err)
		return
	}

	var smtpHost, smtpUsername, smtpPassword, smtpSender, smtpEncryption, adminEmail string
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
		case "admin_email":
			adminEmail = s.Value
		}
	}

	if smtpHost == "" {
		m.logger.Debug("SMTP not configured, skipping email alert")
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

	// Fetch monitor owner
	var user userEntity.User
	if err := m.db.First(&user, incident.UserID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor owner user: %v", err)
		return
	}

	// Send to owner user email
	if user.Email != "" {
		m.logger.Info("Sending email alert to owner: %s", user.Email)
		err := email.SendEmail(cfg, user.Email, subject, body)
		if err != nil {
			m.logger.Error("Failed to send email alert to owner %s: %v", user.Email, err)
		}
	}

	// Send to backup admin email if it's set and different
	if adminEmail != "" && adminEmail != user.Email {
		m.logger.Info("Sending email alert to admin backup: %s", adminEmail)
		err := email.SendEmail(cfg, adminEmail, subject, body)
		if err != nil {
			m.logger.Error("Failed to send email alert to admin backup %s: %v", adminEmail, err)
		}
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

	var discordBotToken, telegramBotToken string
	for _, s := range settings {
		if s.Key == "discord_bot_token" {
			discordBotToken = s.Value
		} else if s.Key == "telegram_bot_token" {
			telegramBotToken = s.Value
		}
	}

	for _, ch := range channels {
		switch ch.Type {
		case "discord_bot":
			if discordBotToken == "" {
				m.logger.Warn("Discord bot token not configured globally, skipping discord_bot channel")
				continue
			}
			go m.sendDiscordBotAlert(ch, discordBotToken, mon, incident, isDown)
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
		case "email":
			go m.sendEmailChannelAlert(ch, mon, incident, isDown)
		}
	}
}

func (m *Module) sendDiscordBotAlert(ch entity.NotificationChannel, token string, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
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

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		m.logger.Error("Failed to create Discord session: %v", err)
		return
	}
	defer dg.Close()

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

	_, err = dg.ChannelMessageSend(config.ChannelID, content)
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

	payload := map[string]interface{}{
		"event": "monitor.status_changed",
		"monitor": map[string]interface{}{
			"id":   mon.ID,
			"name": mon.Name,
			"url":  mon.URL,
		},
		"incident": map[string]interface{}{
			"id":            incident.ID,
			"status":        status,
			"error_message": incident.ErrorMessage,
			"latency":       incident.Latency,
			"created_at":    incident.CreatedAt,
		},
	}
	m.sendJSONPost(config.WebhookURL, payload)
}

func (m *Module) sendEmailChannelAlert(ch entity.NotificationChannel, mon *monitorEntity.Monitor, incident *incidentEntity.Incident, isDown bool) {
	// Fetch monitor owner
	var user userEntity.User
	if err := m.db.First(&user, incident.UserID).Error; err != nil {
		m.logger.Error("Failed to fetch monitor owner user for email channel: %v", err)
		return
	}

	if user.Email == "" {
		m.logger.Warn("Owner email is empty, skipping email channel alert")
		return
	}

	var settings []settingEntity.Setting
	if err := m.db.Find(&settings).Error; err != nil {
		m.logger.Error("Failed to fetch settings for email channel alert: %v", err)
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
		m.logger.Debug("SMTP not configured, skipping email channel alert")
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

	subject := fmt.Sprintf("[ALERT] Monitor %s is %s", mon.Name, map[bool]string{true: "DOWN", false: "UP"}[isDown])
	body := fmt.Sprintf(`
		<h3>Monitor Alert: %s</h3>
		<p><strong>Monitor Name:</strong> %s</p>
		<p><strong>Target URL:</strong> <a href="%s">%s</a></p>
		<p><strong>Status:</strong> <span style="color: %s; font-weight: bold;">%s</span></p>
		<p><strong>Triggered/Resolved At:</strong> %s</p>
		<p><strong>Response Latency:</strong> %dms</p>
	`, map[bool]string{true: "DOWN", false: "UP (RESOLVED)"}[isDown],
		mon.Name, mon.URL, mon.URL,
		map[bool]string{true: "red", false: "green"}[isDown],
		map[bool]string{true: "DOWN", false: "UP"}[isDown],
		incident.CreatedAt.Format("2006-01-02 15:04:05 MST"), incident.Latency)

	err := email.SendEmail(cfg, user.Email, subject, body)
	if err != nil {
		m.logger.Error("Failed to send email channel alert to %s: %v", user.Email, err)
	} else {
		m.logger.Info("Email channel alert sent successfully to %s", user.Email)
	}
}

func (m *Module) sendJSONPost(url string, payload interface{}) {
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

