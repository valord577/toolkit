package email

import (
	"github.com/valord577/mailx"

	"toolkit/system"
)

const (
	mailAlertRecv = "TOOLKIT_MAIL_ALERT_RECV"
)

func Alert(subject, message string) (err error) {
	revc := system.GetEnvString(mailAlertRecv)
	if len(revc) < 1 {
		return
	}

	m := mailx.NewMessage()
	m.SetTo(revc)
	m.SetSubject(subject)
	m.SetPlainBody(message)
	return Send(m)
}
