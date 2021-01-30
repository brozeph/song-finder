// Package services provides utilities for analyzing
// images of songs for the purposes of finding the
// Spotify URI
package services

import (
	"context"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

// AnalyzeImage accepts an image path, reads the image and
// requests to retrieve text annotations from the Google Cloud
// vision API
func AnalyzeImage(imagePath string) (string, error) {
	annotation := ""
	ctx := context.Background()

	f, err := os.Open(imagePath)
	if err != nil {
		return annotation, err
	}

	ifr, err := vision.NewImageFromReader(f)
	if err != nil {
		return annotation, err
	}

	// for each image, upload to Google image analysis
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return annotation, err
	}

	annotations, err := client.DetectTexts(ctx, ifr, nil, 10)
	if err != nil {
		return annotation, err
	}

	if annotations[0] != nil {
		annotation = annotations[0].Description
	}

	return annotation, nil
}
