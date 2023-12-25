package email

import (
	"os"
	"strconv"

	"github.com/valord577/mailx"
)

func Send(m *mailx.Message) (err error) {
	d := &mailx.Dialer{
		Host: smtpHost(),
		Port: smtpPort(),

		Username: smtpUser(),
		Password: smtpPass(),

		SSLOnConnect: smtpTls(),
	}
	return d.DialAndSend(m)
}

const (
	mailSmtpHost   = "TOOLKIT_MAIL_SMTP_HOST"
	mailSmtpPort   = "TOOLKIT_MAIL_SMTP_PORT"
	mailSmtpUser   = "TOOLKIT_MAIL_SMTP_USER"
	mailSmtpPass   = "TOOLKIT_MAIL_SMTP_PASS"
	mailSmtpUseTls = "TOOLKIT_MAIL_SMTP_USE_TLS"
)

func smtpHost() string {
	return os.Getenv(mailSmtpHost)
}

func smtpPort() int {
	env := os.Getenv(mailSmtpPort)
	port, _ := strconv.ParseInt(env, 10, 32)
	return int(port)
}

func smtpUser() string {
	return os.Getenv(mailSmtpUser)
}

func smtpPass() string {
	return os.Getenv(mailSmtpPass)
}

func smtpTls() bool {
	env := os.Getenv(mailSmtpUseTls)
	useTls, _ := strconv.ParseBool(env)
	return useTls
}
