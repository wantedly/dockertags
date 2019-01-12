package main

import (
	"fmt"
	"os"
	"strings"
	"flag"
	"regexp"
)

const Usage = `Usage:
  dockertags IMAGENAME ... [IMAGENAME]
`


func main() {
	imagePtr := flag.Bool("i", false, "show image path in output")
	regexPtr := flag.String("r",".+", "regexp tags")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(1)
	}

	var (
		repo  string
		grepo string
		image string
	)
	for i := 0; i <= flag.NArg() - 1; i++ {

		ss := strings.Split(flag.Arg(i), "/")

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
		        match, _ := regexp.MatchString(*regexPtr, tag)
			if match {
				if flag.NArg() >= 2 || *imagePtr {
					fmt.Printf("%s:%s\n", flag.Arg(i), tag)
				} else {
					fmt.Println(tag)
				}
			}
		}
	}
}
