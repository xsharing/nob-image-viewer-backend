package imageviewerapi

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	session := session.Must(session.NewSession())
	svc := s3.New(session)

	pageNum := 0
	memo := make(map[string]bool)

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("FOLDER_STRUCTURE_BUCKET")),
		//MEMO: 全てを取得したいのでPrefixやDelimiterは使わない
	}

	err = svc.ListObjectsV2Pages(
		params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			fmt.Println(page)

			for _, obj := range page.Contents {
				key := fmt.Sprintf("/%s", aws.StringValue(obj.Key))
				last := path.Base(key)
				dir := path.Dir(key)

				if last == ".keep" {
					for { // 親フォルダへ辿りながら再帰する。.keepが抜けてる場合があるので
						_, has := memo[dir]
						if has {
							break
						} else {
							memo[dir] = true
							res = append(res, dir)
							dir = path.Dir(dir)
						}
					}
				}
			}

			return lastPage
		},
	)

	if err != nil {
		panic(err)
	}

	sort.Strings(res)

	return res, err
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
