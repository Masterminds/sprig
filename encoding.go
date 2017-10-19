package sprig

import (
	"encoding/base32"
	"encoding/base64"
	"net/url"
)

func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func base32encode(v string) string {
	return base32.StdEncoding.EncodeToString([]byte(v))
}

func base32decode(v string) string {
	data, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func urlEncode(v string) string {
	return url.QueryEscape(v)
}

func urlDecode(v string) (string, error) {
	return url.QueryUnescape(v)
}
