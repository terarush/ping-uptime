package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Sender     string
	Encryption string // SSL, TLS, None
}

func SendEmail(cfg SMTPConfig, to string, subject string, body string) error {
	// Create SMTP Auth
	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	// Format header and body
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		cfg.Sender, to, subject, body)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	if cfg.Encryption == "SSL" {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         cfg.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return err
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return err
		}
		defer c.Quit()

		if auth != nil {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}

		if err = c.Mail(cfg.Sender); err != nil {
			return err
		}

		if err = c.Rcpt(to); err != nil {
			return err
		}

		w, err := c.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}
	} else if cfg.Encryption == "TLS" {
		c, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer c.Close()

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         cfg.Host,
		}

		if err = c.StartTLS(tlsconfig); err != nil {
			return err
		}

		if auth != nil {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}

		if err = c.Mail(cfg.Sender); err != nil {
			return err
		}

		if err = c.Rcpt(to); err != nil {
			return err
		}

		w, err := c.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return err
		}
		w.Close()
	} else {
		// Send unencrypted
		err := smtp.SendMail(addr, auth, cfg.Sender, []string{to}, []byte(message))
		if err != nil {
			return err
		}
	}

	return nil
}
