package http

import (
	"io"
	goHttp "net/http"
)

func Get(url string) (responseBody string, responseSize int, err error) {
	resp, err := goHttp.Get(url)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(body), len(body), nil

}
