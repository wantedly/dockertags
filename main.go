package main

import (
	"fmt"
	"os"
	"strings"
)

const Usage = `Usage:
  dockertags IMAGENAME
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(1)
	}

	var (
		repo  string
		image string
		grepo string
	)

	ss := strings.Split(os.Args[1], "/")

	if len(ss) > 2 {
		if strings.Contains(ss[0], "gcr.io") {
			repo = "gcr.io"
			grepo = ss[0]
		} else {
			repo = ss[0]
		}
		image = strings.Join(ss[1:], "/")
	} else if strings.Contains(ss[0], "gcr.io") {
		repo = "gcr.io"
		grepo = ss[0]
		image = strings.Join(ss[1:], "/")
	} else if len(ss) == 1 {
		// Official image of DockerHub
		repo = "hub.docker.com"
		image = strings.Join(append([]string{"library"}, ss...), "/")
	} else {
		repo = "hub.docker.com"
		image = strings.Join(ss, "/")
	}

	var tags []string

	switch repo {
	case "hub.docker.com":
		t, err := retrieveFromDockerHub(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	case "quay.io":
		t, err := retriveFromQuay(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	case "gcr.io":
		t, err := retriveFromgcr(image, grepo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	default:
		fmt.Fprintf(os.Stderr, "Unsupported image repository: %s\n", repo)
		os.Exit(1)
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
}
