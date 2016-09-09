package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const Usage = `Usage:
  dockertags IMAGENAME
`

const QuayURLBase = "https://quay.io/api/v1/repository/"

func constructURL(base, image string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, image, "tag")

	return u.String(), nil
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
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

func retriveFromQuay(image string) ([]string, error) {
	url, err := constructURL(QuayURLBase, image)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	body, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	// TODO: parse JSON
	fmt.Println(body)

	return []string{}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(1)
	}

	var (
		repo  string
		image string
	)

	ss := strings.Split(os.Args[1], "/")

	if len(ss) > 2 {
		repo = ss[0]
		image = strings.Join(ss[1:], "/")
	} else {
		repo = "index.docker.io"
		image = strings.Join(ss, "/")
	}

	var tags []string

	switch repo {
	case "index.docker.io":
	case "quay.io":
		t, err := retriveFromQuay(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	default:

	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
}
