package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
)

const gcrURLBase = "https://"

type gcrTag struct {
	Revision      bool   `json:"revision"`
	StartTs       int    `json:"start_ts"`
	Name          string `json:"name"`
	DockerImageID string `json:"docker_image_id"`
}

type gcrTagsResponse struct {
        Name string   `json:"name"`
        Tags []string `json:"tags"`
}

func constructgcrURL(image string, grepo string) (string, error) {
	u, err := url.Parse(gcrURLBase)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, grepo, "v2", image, "tags/list")
	return u.String(), nil
}

func retriveFromgcr(image string, grepo string) ([]string, error) {
	url, err := constructgcrURL(image, grepo)
	if err != nil {
		return nil, err
	}
	body, err := httpGet(url, os.Getenv("GCR_TOKEN"), nil)
	if err != nil {
		return nil, err
	}
	var resp gcrTagsResponse

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
