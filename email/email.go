package email

import (
	"toolkit/system"

	"github.com/valord577/mailx"
)

const (
	mailSmtpHost   = "TOOLKIT_MAIL_SMTP_HOST"
	mailSmtpPort   = "TOOLKIT_MAIL_SMTP_PORT"
	mailSmtpUser   = "TOOLKIT_MAIL_SMTP_USER"
	mailSmtpPass   = "TOOLKIT_MAIL_SMTP_PASS"
	mailSmtpUseTls = "TOOLKIT_MAIL_SMTP_USE_TLS"
)

func Send(m *mailx.Message) (err error) {
	if m == nil {
		return
	}

	host := system.GetEnvString(mailSmtpHost)
	port := system.GetEnvInt(mailSmtpPort)
	user := system.GetEnvString(mailSmtpUser)
	pass := system.GetEnvString(mailSmtpPass)
	if len(host) < 1 || port < 1 || len(user) < 1 || len(pass) < 1 {
		return
	}

	d := &mailx.Dialer{
		Host: host,
		Port: port,

		Username: user,
		Password: pass,

		SSLOnConnect: system.GetEnvBool(mailSmtpUseTls),
	}
	return d.DialAndSend(m)
}
