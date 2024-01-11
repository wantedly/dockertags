package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// BasicAuthInfo struct mapping user/pass for base auth
type BasicAuthInfo struct {
	Username string
	Password string
}

func httpGet(url, apiToken string, basicAuth *BasicAuthInfo) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+apiToken)
	}

	if basicAuth != nil {
		req.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer bodyClose(resp.Body) // ignoring error

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error: URL: %s, status code: %d,\nbody:\n%s", url, resp.StatusCode, string(b))
	}

	return string(b), nil
}

func bodyClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
