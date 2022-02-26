package imageviewerapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/btcsuite/btcutil/base58"
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

	session := session.Must(session.NewSession())
	svc := s3.New(session)

	pageNum := 0

	hashPrefix := fmt.Sprintf("%s/", hash(*p.Folder))

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("THUMBNAIL_BUCKET")),
		Prefix: aws.String(path.Join(hashPrefix, *p.Folder)),
	}

	err = svc.ListObjectsV2Pages(
		params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			fmt.Println(page)

			for _, obj := range page.Contents {
				thumbailRequest, _ := svc.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(os.Getenv("THUMBNAIL_BUCKET")),
					Key:    obj.Key,
				})
				thumbnailPresigned, _ := thumbailRequest.Presign(1 * time.Hour)
				originalRequest, _ := svc.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(os.Getenv("IMAGE_BUCKET")),
					Key:    aws.String(strings.TrimPrefix(aws.StringValue(obj.Key), hashPrefix)),
				})
				originalPresigned, _ := originalRequest.Presign(1 * time.Hour)

				holder := &imageviewer.Image{
					Name:         aws.String(path.Base(*obj.Key)),
					ThumbnailURL: aws.String(thumbnailPresigned),
					OriginalURL:  aws.String(originalPresigned),
				}

				res = append(res, holder)
			}

			return lastPage
		},
	)

	if err != nil {
		panic(err)
	}

	return res, err
}

func hash(s string) string {
	fmt.Printf("hash:%s\n", s)
	var buffer bytes.Buffer
	if !strings.HasPrefix(s, "/") {
		buffer.WriteString("/")
	}

	buffer.WriteString(s)

	if !strings.HasSuffix(s, "/") {
		buffer.WriteString("/")
	}

	hash := sha256.Sum256([]byte(buffer.String()))
	return base58.Encode(hash[:])
}
