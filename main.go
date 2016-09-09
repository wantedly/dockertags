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
		fmt.Fprintln(os.Stderr, "Retrive from Docker Hub is NOT implemented yet...")
		os.Exit(1)
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
