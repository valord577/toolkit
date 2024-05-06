package autoip

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"toolkit/system"
)

const (
	envLanPrefix = "TOOLKIT_AUTOIP_LAN_PREFIX"
)

func GetLanIp() (s string, err error) {
	prefix := system.GetEnvString(envLanPrefix)
	if len(prefix) < 1 {
		err = errors.New("blank lan addr prefix")
		return
	}

	var addrs []net.Addr
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	var p netip.Prefix
	for _, addr := range addrs {
		ip := addr.String()
		if !strings.HasPrefix(ip, prefix) {
			continue
		}
		if p, err = netip.ParsePrefix(ip); err == nil {
			break
		}
	}
	s = p.Addr().String()
	return
}

const (
	envWanPreferV6 = "TOOLKIT_AUTOIP_WAN_PREFER_V6"
	envWanGetipUrl = "TOOLKIT_AUTOIP_WAN_GETIP_URL"
)

func prefer(network string) string {
	v6 := system.GetEnvBool(envWanPreferV6)
	f := func(prefix string) string {
		s := prefix + "4"
		if v6 {
			s = prefix + "6"
		}
		return s
	}

	if strings.HasPrefix(network, "tcp") {
		return f("tcp")
	}
	if strings.HasPrefix(network, "udp") {
		return f("udp")
	}
	if strings.HasPrefix(network, "ip") {
		return f("ip")
	}
	return network
}

func client() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			DisableCompression:    true,
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          1,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,

			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := &net.Dialer{Timeout: 10 * time.Second, FallbackDelay: -1 * time.Second}
				return dialer.DialContext(ctx, prefer(network), addr)
			},
		},
	}
}

func GetWanIp() (s string, err error) {
	var u *url.URL
	if u, err = url.Parse(system.GetEnvString(envWanGetipUrl)); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = client().Do(
		&http.Request{Method: http.MethodGet, URL: u, Close: true},
	); err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	var payload []byte
	if payload, err = io.ReadAll(resp.Body); err != nil {
		return
	} else {
		payload = bytes.TrimSpace(payload)
	}
	s = string(payload)
	return
}
