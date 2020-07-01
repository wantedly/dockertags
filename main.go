package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/distribution/reference"
)

const Usage = `Usage:
  dockertags IMAGENAME
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(1)
	}

	ref, err := reference.ParseNormalizedNamed(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	repo, image := reference.SplitHostname(ref)

	var tags []string

	switch {
	case repo == "docker.io" || repo == "hub.docker.com":
		t, err := retrieveFromDockerHub(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	case repo == "quay.io":
		t, err := retriveFromQuay(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	case strings.HasSuffix(repo, "amazonaws.com"):
		t, err := retrieveFromECR(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	case strings.Contains(repo, "gcr.io"):
		t, err := retrieveFromGCR(repo, image)
		if err != nil {
			fmt.Println(err)
		}

		tags = t
	case strings.Contains(repo, "gcr.io"):
		t, err := retrieveFromGCR(repo, image)
		if err != nil {
			fmt.Println(err)
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
