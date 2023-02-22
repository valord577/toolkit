package email

import (
	"os"
	"strconv"

	"github.com/valord577/mailx"
)

const (
	toolkitMailSmtpHost   = "TOOLKIT_MAIL_SMTP_HOST"
	toolkitMailSmtpPort   = "TOOLKIT_MAIL_SMTP_PORT"
	toolkitMailSmtpUser   = "TOOLKIT_MAIL_SMTP_USER"
	toolkitMailSmtpPass   = "TOOLKIT_MAIL_SMTP_PASS"
	toolkitMailSmtpUseTls = "TOOLKIT_MAIL_SMTP_USE_TLS"
)

func smtpHost() string {
	return os.Getenv(toolkitMailSmtpHost)
}

func smtpPort() int {
	env := os.Getenv(toolkitMailSmtpPort)
	port, _ := strconv.ParseInt(env, 10, 32)
	return int(port)
}

func smtpUser() string {
	return os.Getenv(toolkitMailSmtpUser)
}

func smtpPass() string {
	return os.Getenv(toolkitMailSmtpPass)
}

func smtpTls() bool {
	env := os.Getenv(toolkitMailSmtpUseTls)
	useTls, _ := strconv.ParseBool(env)
	return useTls
}

func Send(m *mailx.Message) (err error) {
	d := &mailx.Dialer{
		Host: smtpHost(),
		Port: smtpPort(),

		Username: smtpUser(),
		Password: smtpPass(),

		SSLOnConnect: smtpTls(),
	}

	var ser *mailx.Sender
	ser, err = d.Dial()
	if err != nil {
		return err
	}
	defer func() {
		_ = ser.Close()
	}()
	err = ser.Send(m)
	return
}
