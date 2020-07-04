package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sort"

	"golang.org/x/oauth2/google"
)

type gcrTagsResponse struct {
	Manifest map[string]gcrImageDetail `json:"manifest"`
}

// for sorting setups
type gcrImageDetail struct {
	TimeUploadedMs string   `json:"timeUploadedMs"`
	Tag            []string `json:"tag"`
}

type gcrImages []gcrImageDetail

func fetchBearer(repo string, image string) (string, error) {
	credential, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		return "", err
	}

	token, err := credential.TokenSource.Token()
	if err != nil {
		return "", err
	}

	url := constructGCRAuthURL(repo, image)
	body, err := httpGet(url, "", &BasicAuthInfo{
		Username: "_token",
		Password: token.AccessToken,
	})
	if err != nil {
		return "", err
	}

	var resp DockerHubAuthResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return "", err
	}

	return resp.Token, nil
}

func constructGCRAuthURL(repo string, image string) string {
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

func parseGCRTagsResponse(manifests gcrTagsResponse) gcrImages {
	gcrImages := gcrImages{}
	for _, manifest := range manifests.Manifest {
		gcrImages = append(gcrImages, manifest)
	}
	return gcrImages
}

func extractGCRTagNames(images gcrImages) []string {
	tags := []string{}
	sort.Slice(images, func(i, j int) bool {
		return images[i].TimeUploadedMs > images[j].TimeUploadedMs
	}) // sort Newset -> Oldest

	for _, image := range images {
		for _, tag := range image.Tag {
			tags = append(tags, tag)
		}
	}
	return tags
}

func retrieveFromGCR(repo string, image string) ([]string, error) {
	bearer, _ := fetchBearer(repo, image) // continue if fetchBearer fails for public images

	url := constructGCRAPIURL(repo, image)

	body, err := httpGet(url, bearer, nil)
	if err != nil {
		return nil, err
	}

	var resp gcrTagsResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, err
	}

	images := parseGCRTagsResponse(resp)
	tags := extractGCRTagNames(images)
	return tags, nil
}
