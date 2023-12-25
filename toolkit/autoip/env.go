package main

import (
	"os"
	"strings"
)

const autoipAlertRecv = "TOOLKIT_AUTOIP_ALERT_REVC"

func alertRecv() string {
	return getTrimmedEnv(autoipAlertRecv)
}

const autoipGetIpUrl = "TOOLKIT_AUTOIP_PUBLIC_IP_URL"

func getIpUrl() []string {
	urls := getTrimmedEnv(autoipGetIpUrl)
	return strings.Split(urls, " ")
}

const (
	autoipAliEndpoint  = "TOOLKIT_AUTOIP_ALIYUN_ENDPOINT"
	autoipAliAccessKey = "TOOLKIT_AUTOIP_ALIYUN_ACCESS_KEY"
	autoipAliSecretKey = "TOOLKIT_AUTOIP_ALIYUN_SECRET_KEY"
	autoipAliRecordID  = "TOOLKIT_AUTOIP_ALIYUN_RECORD_ID"
)

func endpoint() string {
	endpoint := getTrimmedEnv(autoipAliEndpoint)
	if len(endpoint) < 1 {
		endpoint = "alidns.cn-hangzhou.aliyuncs.com"
	}
	return endpoint
}
func aliAk() string {
	return getTrimmedEnv(autoipAliAccessKey)
}
func aliSk() string {
	return getTrimmedEnv(autoipAliSecretKey)
}
func recid() string {
	return getTrimmedEnv(autoipAliRecordID)
}

func getTrimmedEnv(k string) string {
	return strings.TrimSpace(os.Getenv(k))
}
