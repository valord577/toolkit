package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/valord577/mailx"

	"toolkit/email"
	log "toolkit/logger"
	"toolkit/system"
)

const (
	toolkitNetdevRecv = "TOOLKIT_NETDEV_REVC"
)

func receiver() string {
	return os.Getenv(toolkitNetdevRecv)
}

func notifyInterfaces() (err error) {
	var interfaces []net.Interface
	if interfaces, err = net.Interfaces(); err != nil {
		return
	}

	var l int
	if l = len(interfaces); l < 1 {
		err = errors.New("empty network interface")
		return
	}
	log.Infof("get the system's network interfaces, len: %d", l)

	// collect infos
	b := &strings.Builder{}
	for _, ifi := range interfaces {
		netAddrs, e := ifi.Addrs()
		if e != nil {
			log.Warnf("get network interface address, name: %s, err: %s", ifi.Name, e.Error())
			continue
		}

		fmt.Fprintf(b, "%d: %s: <%s> mtu %d\n", ifi.Index, ifi.Name, ifi.Flags.String(), ifi.MTU)
		fmt.Fprintf(b, "    hardware address: %s\n", ifi.HardwareAddr.String())

		b.WriteString("    network address: ")
		for _, addr := range netAddrs {
			b.WriteString(addr.String())
			b.WriteString(" ")
		}
		b.WriteString("\n")
	}

	// send an email
	hostname := system.Hostname()
	subject := "the system's network interfaces - <" + hostname + ">"
	plainBody := fmt.Sprintf("<%s> - network interfaces shown as below:\n\n%s\n", hostname, b.String())

	m := mailx.NewMessage()
	m.SetTo(receiver())
	m.SetSubject(subject)
	m.SetPlainBody(plainBody)
	log.Debugf("send an email")
	return email.Send(m)
}
