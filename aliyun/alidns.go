package aliyun

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
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

		errmsg := action + ", Status: " + r.Status + ", Response: " + string(bs)
		if r.StatusCode != http.StatusOK {
			return errors.New(errmsg)
		}
		slog.Debug(errmsg)

		if f != nil {
			err = f(bs)
		}
		return err
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
