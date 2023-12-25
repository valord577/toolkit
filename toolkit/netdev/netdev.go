package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/valord577/mailx"

	"toolkit/email"
	"toolkit/logs"
	"toolkit/system"
)

const (
	netdevRecv = "TOOLKIT_NETDEV_REVC"
)

func receiver() string {
	return strings.TrimSpace(os.Getenv(netdevRecv))
}

func netdev() (err error) {
	var interfaces []net.Interface
	if interfaces, err = net.Interfaces(); err != nil {
		return
	}

	// send an email
	hostname := system.Hostname()
	subject := "the system's network interfaces - <" + hostname + ">"
	plainBody := fmt.Sprintf(
		"<%s> - network interfaces shown as below:\n\n%s\n", hostname, rawInterfaces(interfaces),
	)

	m := mailx.NewMessage()
	m.SetTo(receiver())
	m.SetSubject(subject)
	m.SetPlainBody(plainBody)
	return email.Send(m)
}

func rawInterfaces(interfaces []net.Interface) string {
	var l int
	if l = len(interfaces); l < 1 {
		return "empty network interface"
	}
	logs.Infof("get the system's network interfaces, len: %d", l)

	// collect infos
	b := &strings.Builder{}
	for _, ifi := range interfaces {
		fmt.Fprintf(b, "%d: %s: <%s> mtu %d\n", ifi.Index, ifi.Name, ifi.Flags.String(), ifi.MTU)
		fmt.Fprintf(b, "    hardware address: %s\n", ifi.HardwareAddr.String())

		b.WriteString("    network address: ")
		netAddrs, e := ifi.Addrs()
		if e != nil {
			b.WriteString("[err: " + e.Error() + "]")
		} else {
			for _, addr := range netAddrs {
				b.WriteString(addr.String())
				b.WriteString("\n                     ")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
