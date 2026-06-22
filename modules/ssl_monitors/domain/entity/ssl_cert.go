package entity

import "time"

type SSLCert struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	MonitorID     uint      `gorm:"uniqueIndex;not null" json:"monitor_id"`
	Domain        string    `gorm:"type:varchar(255)" json:"domain"`
	Issuer        string    `gorm:"type:varchar(255)" json:"issuer"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidUntil    time.Time `json:"valid_until"`
	DaysRemaining int       `json:"days_remaining"`
	Status        string    `gorm:"type:varchar(20);default:valid" json:"status"`
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CheckedAt     time.Time `json:"checked_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (*SSLCert) TableName() string { return "ssl_certs" }

func NewSSLCert(monitorID uint, domain, issuer string, validFrom, validUntil time.Time, daysRemaining int, status, errMsg string) *SSLCert {
	now := time.Now()
	return &SSLCert{
		MonitorID:     monitorID,
		Domain:        domain,
		Issuer:        issuer,
		ValidFrom:     validFrom,
		ValidUntil:    validUntil,
		DaysRemaining: daysRemaining,
		Status:        status,
		ErrorMessage:  errMsg,
		CheckedAt:     now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
