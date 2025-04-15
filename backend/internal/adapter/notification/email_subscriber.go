package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/rs/zerolog"
)

// EmailConfig contains configuration for email notifications
type EmailConfig struct {
	Enabled      bool
	SMTPServer   string
	SMTPPort     int
	Username     string
	Password     string
	FromAddress  string
	ToAddresses  []string
	MinLevel     AlertLevel
	SubjectPrefix string
}

// EmailSubscriber implements the AlertSubscriber interface for email notifications
type EmailSubscriber struct {
	config EmailConfig
	logger *zerolog.Logger
}

// NewEmailSubscriber creates a new email subscriber
func NewEmailSubscriber(config EmailConfig, logger *zerolog.Logger) *EmailSubscriber {
	return &EmailSubscriber{
		config: config,
		logger: logger,
	}
}

// HandleAlert processes an alert and sends an email if needed
func (s *EmailSubscriber) HandleAlert(alert Alert) error {
	if !s.config.Enabled {
		return nil
	}

	// Check if alert level meets minimum threshold
	if !s.shouldSendAlert(alert) {
		return nil
	}

	// Prepare email content
	subject := s.formatSubject(alert)
	body := s.formatBody(alert)

	// Send email
	return s.sendEmail(subject, body)
}

// GetName returns the name of the subscriber
func (s *EmailSubscriber) GetName() string {
	return "email"
}

// shouldSendAlert determines if an alert should be sent
func (s *EmailSubscriber) shouldSendAlert(alert Alert) bool {
	// Don't send for resolved alerts unless they're critical
	if alert.Resolved && alert.Level != AlertLevelCritical {
		return false
	}

	// Check minimum level
	switch s.config.MinLevel {
	case AlertLevelCritical:
		return alert.Level == AlertLevelCritical
	case AlertLevelError:
		return alert.Level == AlertLevelCritical || alert.Level == AlertLevelError
	case AlertLevelWarning:
		return alert.Level == AlertLevelCritical || alert.Level == AlertLevelError || alert.Level == AlertLevelWarning
	default:
		return true
	}
}

// formatSubject formats the email subject
func (s *EmailSubscriber) formatSubject(alert Alert) string {
	prefix := s.config.SubjectPrefix
	if prefix == "" {
		prefix = "[CryptoBot]"
	}

	var status string
	if alert.Resolved {
		status = "RESOLVED"
	} else {
		status = string(alert.Level)
	}

	return fmt.Sprintf("%s %s: %s", prefix, status, alert.Title)
}

// formatBody formats the email body
func (s *EmailSubscriber) formatBody(alert Alert) string {
	// Simple HTML template for the email
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; color: #333; }
        .alert { border: 1px solid #ddd; border-radius: 5px; padding: 15px; margin-bottom: 20px; }
        .alert-info { border-left: 5px solid #5bc0de; }
        .alert-warning { border-left: 5px solid #f0ad4e; }
        .alert-error { border-left: 5px solid #d9534f; }
        .alert-critical { border-left: 5px solid #d9534f; background-color: #f2dede; }
        .header { font-size: 18px; font-weight: bold; margin-bottom: 10px; }
        .resolved { background-color: #dff0d8; border-color: #d6e9c6; }
        .details { margin-top: 20px; }
        .footer { margin-top: 30px; font-size: 12px; color: #777; }
    </style>
</head>
<body>
    <div class="alert alert-{{.Level}}{{if .Resolved}} resolved{{end}}">
        <div class="header">{{.Title}}</div>
        <p>{{.Message}}</p>
        <div class="details">
            <p><strong>Source:</strong> {{.Source}}</p>
            <p><strong>Time:</strong> {{.Timestamp}}</p>
            <p><strong>Status:</strong> {{if .Resolved}}Resolved{{else}}Active{{end}}</p>
            {{if .Resolved}}<p><strong>Resolved At:</strong> {{.ResolvedAt}}</p>{{end}}
        </div>
    </div>
    <div class="footer">
        <p>This is an automated message from the CryptoBot monitoring system.</p>
    </div>
</body>
</html>
`

	// Parse the template
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to parse email template")
		return fmt.Sprintf("Alert: %s\nMessage: %s\nSource: %s\nTime: %s\nStatus: %s",
			alert.Title, alert.Message, alert.Source, alert.Timestamp.Format(time.RFC3339),
			map[bool]string{true: "Resolved", false: "Active"}[alert.Resolved])
	}

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, alert)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to execute email template")
		return fmt.Sprintf("Alert: %s\nMessage: %s\nSource: %s\nTime: %s\nStatus: %s",
			alert.Title, alert.Message, alert.Source, alert.Timestamp.Format(time.RFC3339),
			map[bool]string{true: "Resolved", false: "Active"}[alert.Resolved])
	}

	return buf.String()
}

// sendEmail sends an email
func (s *EmailSubscriber) sendEmail(subject, body string) error {
	if len(s.config.ToAddresses) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Set up authentication information
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPServer)

	// Prepare email headers
	headers := make(map[string]string)
	headers["From"] = s.config.FromAddress
	headers["To"] = s.config.ToAddresses[0]
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Construct message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Send email
	addr := fmt.Sprintf("%s:%d", s.config.SMTPServer, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, s.config.FromAddress, s.config.ToAddresses, []byte(message))
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to send email")
		return err
	}

	s.logger.Info().
		Str("subject", subject).
		Strs("recipients", s.config.ToAddresses).
		Msg("Email alert sent")

	return nil
}
