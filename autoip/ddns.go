package autoip

import (
	"encoding/json"
	"net/netip"

	"toolkit/aliyun"
	"toolkit/email"
	"toolkit/logs"
	"toolkit/system"
)

const (
	// https://help.aliyun.com/document_detail/2355662.html
	envDdnsAliEndpoint  = "TOOLKIT_DDNS_ALIYUN_ENDPOINT"
	envDdnsAliAccessKey = "TOOLKIT_DDNS_ALIYUN_ACCESS_KEY"
	envDdnsAliSecretKey = "TOOLKIT_DDNS_ALIYUN_SECRET_KEY"
)

func DynamicDNS(ip, recid string) (err error) {
	var addr netip.Addr
	if addr, err = netip.ParseAddr(ip); err != nil {
		return
	}

	var (
		record string
		rtype  string
		rvalue string
	)
	if addr.Is4() {
		rtype = "A"
	}
	if addr.Is6() {
		rtype = "AAAA"
	}

	alidns := aliyun.AliDNS(
		system.GetEnvString(envDdnsAliAccessKey),
		system.GetEnvString(envDdnsAliSecretKey),
	)
	endpoint := system.GetEnvString(envDdnsAliEndpoint)
	err = alidns.DescribeDomainRecordInfo(endpoint, recid,
		func(body []byte) error {
			var v struct {
				RR    string `json:"RR"`
				Value string `json:"Value"`
			}
			e := json.Unmarshal(body, &v)
			if e == nil {
				record = v.RR
				rvalue = v.Value
			}
			return e
		},
	)
	if err != nil || ip == rvalue {
		return
	}
	logs.Infof(
		"ip: '%s', rr: '%s', type: '%s', value: '%s'",
		ip, record, rtype, rvalue,
	)
	subject := "toolkit - autoip <" + system.Hostname() + ">"
	message := "ip: '" + ip + "', rr: '" + record + "', type: '" + rtype + "', value: '" + rvalue + "'"
	if e := email.Alert(subject, message); e != nil {
		logs.Errorf("send alert, err: %s", e.Error())
	}

	return alidns.UpdateDomainRecord(
		endpoint, recid, record, rtype, ip, nil,
	)
}
