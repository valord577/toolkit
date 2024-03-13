package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"toolkit/aliyun/internal"
	"toolkit/logger"

	"golang.org/x/exp/maps"
)

type callback func(*http.Response) error

const (
	RpcFormatXml  = "XML"
	RpcFormatJson = "JSON"
)

var DefaultRpcFormat = RpcFormatJson

func call(
	endpoint string,
	action string,
	version string,
	ak, sk string,
	actionParams map[string]string,
	callback callback,
) error {
	// https://help.aliyun.com/zh/sdk/product-overview/rpc-mechanism
	utc := time.Now().UTC()
	signParams := map[string]string{
		"Action":           action,
		"Version":          version,
		"Format":           DefaultRpcFormat,
		"AccessKeyId":      ak,
		"SignatureNonce":   utc.Format(time.RFC3339Nano),
		"Timestamp":        utc.Format(time.RFC3339),
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
	}
	maps.Copy(signParams, actionParams)

	method := http.MethodPost
	body := getRequestStr(method, sk, signParams)

	request := &http.Request{
		URL: &url.URL{
			Scheme: "https",
			Host:   endpoint,
		},
		Header: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},

		Method:        method,
		ContentLength: int64(len(body)),
		Body:          io.NopCloser(strings.NewReader(body)),
		Close:         true,
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return callback(resp)
}

func getRequestStr(method, sk string, signParams map[string]string) string {
	buf := &strings.Builder{}

	keys := maps.Keys(signParams)
	sort.Strings(keys)
	for _, k := range keys {
		buf.WriteString(internal.Escape(k))
		buf.WriteByte('=')
		buf.WriteString(internal.Escape(signParams[k]))
		buf.WriteByte('&')
	}
	noSignStr := buf.String()[:buf.Len()-1]
	encodeStr := internal.Escape(noSignStr)

	buf.Reset()
	buf.WriteString(method)
	buf.WriteString("&%2F&")
	buf.WriteString(encodeStr)
	strToSign := buf.String()
	logger.Debugf("strToSign: %s", strToSign)

	// hmac-sha1
	w := hmac.New(sha1.New, []byte(sk+"&"))
	_, _ = w.Write([]byte(strToSign))
	sign := base64.StdEncoding.EncodeToString(w.Sum(nil))
	return "Signature=" + sign + "&" + noSignStr
}
