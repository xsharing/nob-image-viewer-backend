// Code generated by goa v3.5.4, DO NOT EDIT.
//
// image-viewer HTTP client CLI support package
//
// Command:
// $ goa gen github.com/image-viewer/nob-image-viewer-backend/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	imageviewerc "github.com/image-viewer/nob-image-viewer-backend/gen/http/image_viewer/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `image-viewer (folders|images)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` image-viewer folders` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		imageViewerFlags = flag.NewFlagSet("image-viewer", flag.ContinueOnError)

		imageViewerFoldersFlags = flag.NewFlagSet("folders", flag.ExitOnError)

		imageViewerImagesFlags      = flag.NewFlagSet("images", flag.ExitOnError)
		imageViewerImagesFolderFlag = imageViewerImagesFlags.String("folder", "", "")
	)
	imageViewerFlags.Usage = imageViewerUsage
	imageViewerFoldersFlags.Usage = imageViewerFoldersUsage
	imageViewerImagesFlags.Usage = imageViewerImagesUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "image-viewer":
			svcf = imageViewerFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "image-viewer":
			switch epn {
			case "folders":
				epf = imageViewerFoldersFlags

			case "images":
				epf = imageViewerImagesFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "image-viewer":
			c := imageviewerc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "folders":
				endpoint = c.Folders()
				data = nil
			case "images":
				endpoint = c.Images()
				data, err = imageviewerc.BuildImagesPayload(*imageViewerImagesFolderFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// image-viewerUsage displays the usage of the image-viewer command and its
// subcommands.
func imageViewerUsage() {
	fmt.Fprintf(os.Stderr, `The image-viewer service
Usage:
    %[1]s [globalflags] image-viewer COMMAND [flags]

COMMAND:
    folders: Folders implements folders.
    images: Images implements images.

Additional help:
    %[1]s image-viewer COMMAND --help
`, os.Args[0])
}
func imageViewerFoldersUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] image-viewer folders

Folders implements folders.

Example:
    %[1]s image-viewer folders
`, os.Args[0])
}

func imageViewerImagesUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] image-viewer images -folder STRING

Images implements images.
    -folder STRING: 

Example:
    %[1]s image-viewer images --folder "Dolores est."
`, os.Args[0])
}
