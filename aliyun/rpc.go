package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"toolkit/system"
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
	for k, v := range actionParams {
		signParams[k] = v
	}

	method := http.MethodPost
	body := getRequestStr(method, sk, signParams)

	request := &http.Request{
		URL: &url.URL{
			Scheme: "https",
			Host:   endpoint,
		},
		Header: http.Header{
			"Content-Type": {"application/x-www-form-urlencoded; charset=utf-8"},
			"User-Agent":   {system.Version()},
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

	keys := make([]string, 0, len(signParams))
	for k := range signParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		buf.WriteString(system.Escape(k))
		buf.WriteByte('=')
		buf.WriteString(system.Escape(signParams[k]))
		buf.WriteByte('&')
	}
	noSignStr := buf.String()[:buf.Len()-1]
	encodeStr := system.Escape(noSignStr)

	buf.Reset()
	buf.WriteString(method)
	buf.WriteString("&%2F&")
	buf.WriteString(encodeStr)
	strToSign := buf.String()
	slog.Debug("strToSign: " + strToSign)

	// hmac-sha1
	w := hmac.New(sha1.New, []byte(sk+"&"))
	_, _ = w.Write([]byte(strToSign))
	sign := base64.StdEncoding.EncodeToString(w.Sum(nil))
	return "Signature=" + sign + "&" + noSignStr
}
