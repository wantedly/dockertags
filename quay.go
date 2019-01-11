package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
	"strconv"
)

const QuayURLBase = "https://quay.io/api/v1/repository/"

type QuayTag struct {
	Revision      bool   `json:"revision"`
	StartTs       int    `json:"start_ts"`
	Name          string `json:"name"`
	DockerImageID string `json:"docker_image_id"`
}

type QuayTagsResponse struct {
	HasAdditional bool      `json:"has_additional"`
	Page          int       `json:"page"`
	Tags          []QuayTag `json:"tags"`
}

func constructQuayURL(image string, Page int) (string, error) {
	u, err := url.Parse(QuayURLBase)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, image, "tag") + "/"
	queryString := u.Query()
	queryString.Set("page", strconv.Itoa(Page))
	u.RawQuery = queryString.Encode()

	return u.String(), nil
}

func retriveFromQuay(image string) ([]string, error) {
	var (
  	HasAdditional	bool
    Page		int
		resp		QuayTagsResponse
  )
  HasAdditional = true
  Page = 1
  tags := []string{}

  for HasAdditional {
		url, err := constructQuayURL(image, Page)
		if err != nil {
			return nil, err
		}
		body, err := httpGet(url, os.Getenv("QUAYIO_TOKEN"), nil)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			return nil, err
		}
		Page = resp.Page+1
		HasAdditional = resp.HasAdditional

		for _, tag := range resp.Tags {
			tags = append(tags, tag.Name)
		}
	}

	return tags, nil
}
