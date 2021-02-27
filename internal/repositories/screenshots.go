package repositories

import (
	"context"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/brozeph/song-finder/internal/interfaces"
	"github.com/brozeph/song-finder/internal/models"
	"github.com/rs/zerolog/log"
)

var imageExtensions = []string{".png", ".jpg", ".jpeg"}

type screenshotRepository struct{}

// NewScreenshotRepository returns a new instance
func NewScreenshotRepository() interfaces.IScreenshotRepository {
	return &screenshotRepository{}
}

// DetectText accepts an image path, reads the image and
// requests to retrieve text annotations from the Google Cloud
// vision API
func (sr *screenshotRepository) DetectText(path string) (string, error) {
	text := ""
	ctx := context.Background()

	f, err := os.Open(path)
	if err != nil {
		return text, err
	}

	ifr, err := vision.NewImageFromReader(f)
	if err != nil {
		return text, err
	}

	// for each image, upload to Google image analysis
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return text, err
	}

	annotations, err := client.DetectTexts(ctx, ifr, nil, 10)
	if err != nil {
		return text, err
	}

	if annotations[0] != nil {
		text = annotations[0].Description
	}

	return text, nil
}

func (sr *screenshotRepository) FindInPath(path string) ([]*models.Screenshot, error) {
	var (
		h  = sha256.New()
		sf []*models.Screenshot
	)

	// find all of the image files
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if isImage(path) && info.Size() > 0 {
			// reset hash
			h.Reset()

			// determine shasum
			f, err := os.Open(path)
			if err != nil {
				return nil
			}

			defer f.Close()
			if _, err := io.Copy(h, f); err != nil {
				log.Error().Stack().Err(err).Msg("unable to calculate sha sum")
			}

			sf = append(sf, &models.Screenshot{
				Path:   path,
				SHASum: string(h.Sum(nil)),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sf, nil
}

func isImage(filePath string) bool {
	var fileExt = strings.ToLower(filepath.Ext(filePath))

	for _, ext := range imageExtensions {
		if fileExt == ext {
			return true
		}
	}

	return false
}
