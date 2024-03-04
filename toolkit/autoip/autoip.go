package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/valord577/mailx"

	"toolkit/aliyun"
	"toolkit/email"
	"toolkit/logs"
	"toolkit/system"
)

func getPublicIpFromHTTP(httpUrl string) (ip string, err error) {
	var u *url.URL
	if u, err = url.Parse(httpUrl); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(
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
	}
	payload = bytes.TrimRightFunc(payload, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	ip = string(payload)
	logs.Infof("ip: '%s', url: '%s'", ip, httpUrl)
	return
}

func getPublicIp() (ip string, err error) {
	for _, u := range getIpUrl() {
		if len(u) < 1 {
			continue
		}
		if ip, err = getPublicIpFromHTTP(u); err != nil {
			continue
		}
		if len(ip) > 0 {
			break
		}
	}
	if len(ip) < 1 {
		if err == nil {
			err = errors.New("failed to get public ip")
		}
	} else {
		err = nil
	}
	return
}

func autoip() (err error) {
	defer func() {
		if err != nil {
			alertErr(err)
		}
	}()

	var (
		ipstr  string
		rtype  string
		record string
		rvalue string
	)
	if ipstr, err = getPublicIp(); err != nil {
		return
	}
	ip := net.ParseIP(ipstr)
	if ip == nil {
		err = errors.New("invalid ip addr: " + ipstr)
		return
	}
	if netip := ip.To4(); netip != nil {
		rtype = "A"
	} else {
		rtype = "AAAA"
	}

	alidns := aliyun.AliDNS(aliAk(), aliSk())
	err = alidns.DescribeDomainRecordInfo(
		endpoint(), recid(),

		func(resp *http.Response) (e error) {
			var v struct {
				RR    string `json:"RR"`
				Value string `json:"Value"`
			}
			e = json.NewDecoder(resp.Body).Decode(&v)
			if e == nil {
				record = v.RR
				rvalue = v.Value
			}
			return e
		},
	)
	if err != nil || ipstr == rvalue {
		return
	}

	logs.Infof(
		"ip: '%s', rr: '%s', type: '%s', value: '%s'",
		ipstr, record, rtype, rvalue,
	)
	return alidns.UpdateDomainRecord(
		endpoint(), recid(), record, rtype, ipstr,

		func(resp *http.Response) (e error) {
			b := &strings.Builder{}
			b.WriteString("\n>>>>>>>>>>\n")
			e = resp.Write(b)
			b.WriteString("\n<<<<<<<<<<\n")

			logs.Infof("UpdateDomainRecord Response: %s", b.String())
			return e
		},
	)
}

func alertErr(err error) {
	recv := alertRecv()
	if len(recv) < 1 {
		return
	}

	m := mailx.NewMessage()
	m.SetTo(recv)
	m.SetSubject("autoip - update alidns domain record - <" + system.Hostname() + ">")
	m.SetPlainBody("error: " + err.Error())
	if e := email.Send(m); e != nil {
		logs.Errorf("send alert email, err: %s", e.Error())
	}
}
