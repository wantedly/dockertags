package main

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func retrieveFromECR(image string) ([]string, error) {
	session, err := session.NewSession()
	allTags := []string{}
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, aerr
		}
		return nil, err
	}

	svc := ecr.New(session)
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(image),
		Filter: &ecr.DescribeImagesFilter{
			TagStatus: aws.String("TAGGED"), // extract tagged images only
		},
	}

  for 	{
   	result, err := svc.DescribeImages(input)
   	if err != nil {
   		if aerr, ok := err.(awserr.Error); ok {
   			return nil, aerr
   		}
   		return nil, err
   	}

   	tags := extractEcrTagNames(result.ImageDetails)
		allTags = append(allTags, tags...)
		// Check if there are more results
		if result.NextToken == nil {
			// No more results, break out of the loop
			break
		}
		// Set the next token for the next iteration
		input.NextToken = result.NextToken
	}
	return allTags	, nil
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
