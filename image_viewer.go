package imageviewerapi

import (
	"context"
	"log"

	imageviewer "github.com/image-viewer/nob-image-viewer-backend/gen/image_viewer"
)

// image-viewer service example implementation.
// The example methods log the requests and return zero values.
type imageViewersrvc struct {
	logger *log.Logger
}

// NewImageViewer returns the image-viewer service implementation.
func NewImageViewer(logger *log.Logger) imageviewer.Service {
	return &imageViewersrvc{logger}
}

// Folders implements folders.
func (s *imageViewersrvc) Folders(ctx context.Context) (res []string, err error) {
	s.logger.Print("imageViewer.folders")
	return
}

// Images implements images.
func (s *imageViewersrvc) Images(ctx context.Context, p *imageviewer.ImagesPayload) (res []*imageviewer.Image, err error) {
	s.logger.Print("imageViewer.images")
	return
}
