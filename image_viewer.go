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
	return []string{"/", "/folder"}, nil
}

func Const2Ptr(str string) *string {
	return &str
}

// Images implements images.
func (s *imageViewersrvc) Images(ctx context.Context, p *imageviewer.ImagesPayload) (res []*imageviewer.Image, err error) {
	s.logger.Print("imageViewer.images")
	return []*imageviewer.Image{
		{Name: Const2Ptr("a"), ThumbnailURL: Const2Ptr("https://dummyimage.com/200x200/000/00f"), OriginalURL: Const2Ptr("https://dummyimage.com/1280x1280/000/00f")},
		{Name: Const2Ptr("b"), ThumbnailURL: Const2Ptr("https://dummyimage.com/200x200/000/00f"), OriginalURL: Const2Ptr("https://dummyimage.com/1280x1280/000/00f")},
	}, nil
}
