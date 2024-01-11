package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/docker/distribution/reference"
)

// Usage string
const Usage = `Usage: dockertags IMAGENAME [REGEXP]`

func main() {
	var tags []string
	var re *regexp.Regexp

	if len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(1)
	}

	ref, err := reference.ParseNormalizedNamed(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	repo, image := reference.SplitHostname(ref)

	if len(os.Args) == 3 {
		re, _ = regexp.Compile(os.Args[2])
	}

	switch {
	case repo == "docker.io" || repo == "hub.docker.com":
		t, err := retrieveFromDockerHub(image)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tags = t
	case repo == "quay.io":
		t, err := retrieveFromQuay(image)
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
	case strings.HasSuffix(repo, "gcr.io"):
		t, err := retrieveFromGCR(repo, image)
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
		if re == nil {
			fmt.Println(tag)
		} else if re.MatchString(tag) {
			fmt.Println(tag)
		}
	}
}
