package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpGet(url, apiToken string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	if apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+apiToken)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error!\nURL: %s\nstatus code: %d\nbody:\n%s\n", url, resp.StatusCode, string(b))
	}

	return string(b), nil
}
