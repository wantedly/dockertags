package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func httpGet(url string) (string, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("QUAYIO_TOKEN"))
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
