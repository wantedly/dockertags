package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// for sorting setups
// TODO: consider using global parameter 'query' to get sorted result directory from API
// ref: https://docs.aws.amazon.com/cli/latest/reference/index.html
type Images []ecr.ImageDetail // type alias to implement Len and Swap

func (img Images) Len() int {
	return len(img)
}

func (img Images) Swap(i, j int) {
	img[i], img[j] = img[j], img[i]
}

type ByPushedAt struct {
	Images
}

func (b ByPushedAt) Less(i, j int) bool {
	ti := *b.Images[i].ImagePushedAt
	tj := *b.Images[j].ImagePushedAt
	return ti.Before(tj)
}

func retrieveFromECR(image string) ([]string, error) {
	profile := os.Getenv("AWS_PROFILE")

	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("Error: AWS_REGION must be set")
	}

	svc := ecr.New(session.New(), &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(image),
		Filter: &ecr.DescribeImagesFilter{
			TagStatus: aws.String("TAGGED"), // extract tagged images only
		},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf(aerr.Error())
		}
		return nil, err
	}

	tags := extractEcrTagNames(castEcrImageDetailsToImages(result.ImageDetails))
	return tags, nil
}

func extractEcrTagNames(images Images) []string {
	tags := []string{}
	sort.Slice(images, func(i, j int) bool {
		return images[i].ImagePushedAt.After(*images[j].ImagePushedAt)
	}) // sort Newest -> Oldest

	for _, image := range images {
		for _, tag := range image.ImageTags {
			tags = append(tags, *tag)
		}
	}
	return tags
}

func castEcrImageDetailsToImages(ecrImages []*ecr.ImageDetail) Images {
	images := Images{}
	for _, image := range ecrImages {
		images = append(images, *image)
	}
	return images
}
