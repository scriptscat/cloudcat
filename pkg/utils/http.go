package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpPost(url, data string, header http.Header) ([]byte, error) {
	method := "POST"

	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	if header == nil {
		req.Header = http.Header{}
	} else {
		req.Header = header.Clone()
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
