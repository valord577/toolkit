package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/netip"

	"toolkit/aliyun"
	"toolkit/logger"
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
		func(resp *http.Response) (e error) {
			var body []byte
			if body, e = io.ReadAll(resp.Body); e != nil {
				return
			}
			logger.Infof("DescribeDomainRecordInfo, Status: %s, Response: %s", resp.Status, body)

			var v struct {
				RR    string `json:"RR"`
				Value string `json:"Value"`
			}
			e = json.Unmarshal(body, &v)
			if e == nil {
				record = v.RR
				rvalue = v.Value
			}
			return
		},
	)
	if err != nil || ip == rvalue {
		return
	}
	logger.Infof(
		"ip: '%s', rr: '%s', type: '%s', value: '%s'",
		ip, record, rtype, rvalue,
	)
	return alidns.UpdateDomainRecord(
		endpoint, recid, record, rtype, ip,

		func(resp *http.Response) (e error) {
			var body []byte
			if body, e = io.ReadAll(resp.Body); e != nil {
				return
			}
			logger.Infof("UpdateDomainRecord, Status: %s, Response: %s", resp.Status, body)
			return
		},
	)
}
