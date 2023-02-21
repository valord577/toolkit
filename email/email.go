package email

import (
	"toolkit/config"
	log "toolkit/logger"

	"github.com/valord577/mailx"
)

func Send(m *mailx.Message) {
	cnf := config.SmtpCnf()

	d := &mailx.Dialer{
		Host: cnf.Host,
		Port: cnf.Port,

		Username: cnf.User,
		Password: cnf.Pass,

		SSLOnConnect: cnf.SslOnConnect,
	}

	var (
		err error
		ser *mailx.Sender
	)
	if ser, err = d.Dial(); err != nil {
		log.Warnf("can not dial to smtp server, err: %s", err.Error())
		return
	}
	if err = ser.Send(m); err != nil {
		log.Errorf("send email, err: %s", err.Error())
		return
	}
	if err = ser.Close(); err != nil {
		log.Warnf("close smtp pipe, err: %s", err.Error())
	}
}
