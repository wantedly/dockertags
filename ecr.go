package main

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func retrieveFromECR(image string) ([]string, error) {
	svc := ecr.New(session.New())
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

	tags := extractEcrTagNames(result.ImageDetails)
	return tags, nil
}

func extractEcrTagNames(images []*ecr.ImageDetail) []string {
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
