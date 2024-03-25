package aliyun

import (
	"fmt"
	"io"
	"net/http"

	"toolkit/logs"
)

type alidns struct {
	ak string
	sk string
}

func AliDNS(ak, sk string) *alidns {
	return &alidns{ak, sk}
}

func (c *alidns) callback(
	action string, f func(body []byte) error,
) callback {
	return func(r *http.Response) error {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		format := "%s, Status: %s, Response: %s"
		args := []any{action, r.Status, string(bs)}
		if r.StatusCode != http.StatusOK {
			return fmt.Errorf(format, args...)
		}
		logs.Debugf(format, args...)

		if f != nil {
			return f(bs)
		}
		return nil
	}
}

func (c *alidns) DescribeDomainRecordInfo(
	endpoint, rid string, f func(body []byte) error,
) error {
	// https://help.aliyun.com/document_detail/2357158.html
	const (
		action  = "DescribeDomainRecordInfo"
		version = "2015-01-09"
	)
	actionParams := map[string]string{
		"RecordId": rid,
	}
	return call(
		endpoint, action, version,
		c.ak, c.sk, actionParams, c.callback(action, f),
	)
}

func (c *alidns) UpdateDomainRecord(
	endpoint, rid, rr, rtype, rvalue string,
	f func(body []byte) error,
) error {
	// https://help.aliyun.com/document_detail/2355677.html
	const (
		action  = "UpdateDomainRecord"
		version = "2015-01-09"
	)
	actionParams := map[string]string{
		"RecordId": rid,
		"RR":       rr,
		"Type":     rtype,
		"Value":    rvalue,
	}
	return call(
		endpoint, action, version,
		c.ak, c.sk, actionParams, c.callback(action, f),
	)
}
