package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
)

const gcrioURLBase = "https://"

type gcrioTag struct {
	Revision      bool   `json:"revision"`
	StartTs       int    `json:"start_ts"`
	Name          string `json:"name"`
	DockerImageID string `json:"docker_image_id"`
}

type gcrioTagsResponse struct {
        Name string   `json:"name"`
        Tags []string `json:"tags"`
}

func constructgcrioURL(image string, grepo string) (string, error) {
	u, err := url.Parse(gcrioURLBase)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, grepo, "v2", image, "tags/list")
	return u.String(), nil
}

func retriveFromgcrio(image string, grepo string) ([]string, error) {
	url, err := constructgcrioURL(image, grepo)
	if err != nil {
		return nil, err
	}
	body, err := httpGet(url, os.Getenv("GCRIO_TOKEN"), nil)
	if err != nil {
		return nil, err
	}
	var resp gcrioTagsResponse

	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, err
	}
	tags := resp.Tags
  // Reverse the order of the tags to make it ordered as: "latest => oldest"
  for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
          tags[i], tags[j] = tags[j], tags[i]
  }
	return tags, nil
}
