package sprig

import (
	"net/url"
)

func pathEscape(v string) string {
	return url.PathEscape(v)
}

func pathUnescape(v string) string {
	data, err := url.PathUnescape(v)
	if err != nil {
		return err.Error()
	}
	return data
}

func queryEscape(v string) string {
	return url.QueryEscape(v)
}

func queryUnescape(v string) string {
	data, err := url.QueryUnescape(v)
	if err != nil {
		return err.Error()
	}
	return data
}
