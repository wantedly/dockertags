package main

import (
	"fmt"
	"os"

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

	switch repo {
	case "docker.io", "hub.docker.com":
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
	default:
		fmt.Fprintf(os.Stderr, "Unsupported image repository: %s\n", repo)
		os.Exit(1)
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
}
