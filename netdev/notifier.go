package netdev

import (
	"fmt"
	"net"
	"strings"

	"toolkit/email"
	log "toolkit/logger"

	"github.com/valord577/mailx"
)

// interface name -> interface addrs
var cache cachedMapNetAddr

func NotifyChanges(device string, receiver string) {
	nets, err := net.Interfaces()
	if err != nil {
		log.Errorf("get the system's network interfaces, err: %s", err.Error())
		return
	}

	l := len(nets)
	if l < 1 {
		log.Warnf("empty network interface")
		return
	}
	log.Infof("get the system's network interfaces, len: %d", l)

	// collect infos
	current := make(cachedMapNetAddr)
	for _, n := range nets {
		addrs, err := n.Addrs()
		if err != nil {
			log.Errorf("get network address from '%s', err: %s", n.Name, err.Error())
			return
		}
		current.add(n.Name, addrs)
	}
	// diff infos
	if ok := current.eq(cache); ok {
		log.Debugf("skip notification")
		return
	}
	cache = current
	log.Debugf("send notification")
	// nets to string
	sendMail(nets2string(device, nets), device, receiver)
}

func nets2string(device string, nets []net.Interface) string {
	if len(nets) == 0 {
		return ""
	}

	b := &strings.Builder{}
	for _, n := range nets {
		netAddrs, err := n.Addrs()
		if err != nil {
			log.Errorf("get network address from '%s', err: %s", n.Name, err.Error())
			return ""
		}

		fmt.Fprintf(b, "%d: %s: <%s> mtu %d\n", n.Index, n.Name, n.Flags.String(), n.MTU)
		fmt.Fprintf(b, "    hardware address: %s\n", n.HardwareAddr.String())

		b.WriteString("    network address: ")
		for _, addr := range netAddrs {
			b.WriteString(addr.String())
			b.WriteString(" ")
		}
		b.WriteString("\n")
	}
	return fmt.Sprintf("<%s> - network interfaces shown as below:\n\n%s\n", device, b.String())
}

func sendMail(plainBody string, device string, receiver string) {
	subject := "the system's network interfaces - <" + device + ">"

	m := mailx.NewMessage()
	m.SetTo(receiver)
	m.SetSubject(subject)
	m.SetPlainBody(plainBody)
	email.Send(m)
}
