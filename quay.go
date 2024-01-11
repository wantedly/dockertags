package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
)

// QuayURLBase base repository url
const QuayURLBase = "https://quay.io/api/v1/repository/"

// QuayTag struct mapped to Quay
type QuayTag struct {
	Revision      bool   `json:"revision"`
	StartTs       int    `json:"start_ts"`
	Name          string `json:"name"`
	DockerImageID string `json:"docker_image_id"`
}

// QuayTagsResponse list of tags from quay
type QuayTagsResponse struct {
	HasAdditional bool      `json:"has_additional"`
	Page          int       `json:"page"`
	Tags          []QuayTag `json:"tags"`
}

func constructQuayURL(image string) (string, error) {
	u, err := url.Parse(QuayURLBase)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, image, "tag") + "/"

	return u.String(), nil
}

func retrieveFromQuay(image string) ([]string, error) {
	url, err := constructQuayURL(image)
	if err != nil {
		return nil, err
	}

	body, err := httpGet(url, os.Getenv("QUAYIO_TOKEN"), nil)
	if err != nil {
		return nil, err
	}

	var resp QuayTagsResponse

	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, err
	}

	tags := []string{}

	for _, tag := range resp.Tags {
		tags = append(tags, tag.Name)
	}

	return tags, nil
}
