package aliyun

type alidns struct {
	ak string
	sk string
}

func AliDNS(ak, sk string) *alidns {
	return &alidns{ak, sk}
}

func (c *alidns) DescribeDomainRecordInfo(
	endpoint, rid string, callback callback,
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
		c.ak, c.sk, actionParams, callback,
	)
}

func (c *alidns) UpdateDomainRecord(
	endpoint, rid, rr, rtype, rvalue string, callback callback,
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
		c.ak, c.sk, actionParams, callback,
	)
}
