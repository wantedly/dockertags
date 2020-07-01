package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"sort"
	"strings"
)

type gcrAuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type gcrTagsResponse struct {
	Manifest map[string]json.RawMessage `json:"manifest"`
}

// for sorting setups
type gcrImageDetail struct {
	TimeCreatedMs string   `json:"timeCreatedMs"`
	Tag           []string `json:"tag"`
}

type gcrImages []gcrImageDetail

// ref: https://stackoverflow.com/questions/34037256/does-google-container-registry-support-docker-remote-api-v2/34046435#34046435
func fetchBearer(repo string, image string) (string, error) {
	token, err := exec.Command("gcloud", "auth", "print-access-token").Output() // run gcloud command
	if err != nil {
		return "", err
	}

	url := constructAuthURL(repo, image)
	body, err := httpGet(url, "", &BasicAuthInfo{
		Username: "_token",
		Password: strings.TrimSpace(string(token)), // need to remove trailing new line character
	})

	var resp gcrAuthResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return "", err
	}
	return resp.Token, nil
}

func constructAuthURL(repo string, image string) string {
	u := url.URL{
		Scheme: "https",
		Path:   path.Join(repo, "/v2/token"),
	}
	q := u.Query()
	q.Set("scope", fmt.Sprintf("repository:%s:pull", image))
	u.RawQuery = q.Encode()

	return u.String()
}

func constructGCRAPIURL(repo string, image string) string {
	u := url.URL{
		Scheme: "https",
		Path:   path.Join(repo, "v2", image, "tags/list"),
	}
	return u.String()
}

func parseGCRTagsResponse(manifests gcrTagsResponse) (gcrImages, error) {
	gcrImages := gcrImages{}
	for _, manifest := range manifests.Manifest {
		var imageDetail gcrImageDetail
		if err := json.Unmarshal([]byte(manifest), &imageDetail); err != nil {
			return nil, err
		}
		gcrImages = append(gcrImages, imageDetail)
	}
	return gcrImages, nil
}

func extractGCRTagNames(images gcrImages) []string {
	tags := []string{}
	sort.Slice(images, func(i, j int) bool {
		return images[i].TimeCreatedMs > images[j].TimeCreatedMs
	}) // sort Newset -> Oldest

	for _, image := range images {
		for _, tag := range image.Tag {
			tags = append(tags, tag)
		}
	}
	return tags
}

func retrieveFromGCR(repo string, image string) ([]string, error) {
	bearer, err := fetchBearer(repo, image)
	if err != nil {
		return nil, err
	}

	url := constructGCRAPIURL(repo, image)

	body, err := httpGet(url, bearer, nil)
	if err != nil {
		return nil, err
	}

	var resp gcrTagsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, err
	}

	images, err := parseGCRTagsResponse(resp)
	if err != nil {
		return nil, err
	}

	tags := extractGCRTagNames(images)
	return tags, nil
}
