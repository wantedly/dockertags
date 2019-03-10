package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
)

const ecrURLBase = "https://"

func constructEcrURL(image string, grepo string) (string, error) {
	u, err := url.Parse(ecrURLBase)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, grepo, "v2", image, "tags/list")
	return u.String(), nil
}

func retriveFromEcr(image string, grepo string) ([]string, error) {
	url, err := constructEcrURL(image, grepo)
	if err != nil {
		return nil, err
	}
	body, err := httpGet(url, os.Getenv("ECR_TOKEN"), nil)
	if err != nil {
		return nil, err
	}
	var resp ecrTagsResponse

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
