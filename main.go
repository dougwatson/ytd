package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"runtime"

	//"os"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/ioprogress"

	"github.com/dougwatson/youtube/v2"
)

const defaultExtension = ".mov"

// Rely on hardcoded canonical mime types, as the ones provided by Go aren't exhaustive [1].
// This seems to be a recurring problem for youtube downloaders, see [2].
// The implementation is based on mozilla's list [3], IANA [4] and Youtube's support [5].
// [1] https://github.com/golang/go/blob/ed7888aea6021e25b0ea58bcad3f26da2b139432/src/mime/type.go#L60
// [2] https://github.com/ZiTAL/youtube-dl/blob/master/mime.types
// [3] https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
// [4] https://www.iana.org/assignments/media-types/media-types.xhtml#video
// [5] https://support.google.com/youtube/troubleshooter/2888402?hl=en
var canonicals = map[string]string{
	"video/quicktime":  ".mov",
	"video/x-msvideo":  ".avi",
	"video/x-matroska": ".mkv",
	"video/mpeg":       ".mpeg",
	"video/webm":       ".webm",
	"video/3gpp2":      ".3g2",
	"video/x-flv":      ".flv",
	"video/3gpp":       ".3gp",
	"video/mp4":        ".mp4",
	"video/ogg":        ".ogv",
	"video/mp2t":       ".ts",
}

// Downloader offers high level functions to download videos into files
type Downloader struct {
	youtube.Client
	OutputDir string // optional directory to store the files
}

var testDownloader = func() (dl Downloader) {
	dl.OutputDir = "download_test"
	dl.Debug = true
	return
}()

func (dl *Downloader) logf(format string, v ...interface{}) {
	if dl.Debug {
		log.Printf(format, v...)
	}
}

func main() {
	ctx := context.Background()

	video, err := testDownloader.Client.GetVideoContext(ctx, "446E-r0rXHI") //youtube.com
	if err != nil {
		println("HERE")
		panic(err)
	}
	println("main.go main() about to Download ======================================================================")
	//testDownloader.DownloadComposite(ctx, "", video, "hd1080", "mp4")
	printQuality(video)
	testDownloader.Download(ctx, video, &video.Formats[27], "")
	//testDownloader.Download(ctx, video, "hd720", "")
}
func printQuality(arr *youtube.Video) {

	for i, val := range arr.Formats {
		fmt.Printf("quality[%v]=%v\n", i, val.Quality)
		//if list[i].Quality == quality || list[i].QualityLabel == quality {
		//      return &list[i]
		//}
	}
}

// Download : Starting download video by arguments.
func (dl *Downloader) Download(ctx context.Context, v *youtube.Video, format *youtube.Format, outputFile string) error {
	dl.logf("Video '%s' - Quality '%s' - Codec '%s'", v.Title, format.QualityLabel, format.MimeType)
	destFile, err := dl.getOutputFile(v, format, outputFile)
	if err != nil {
		return err
	}

	dl.logf("Download to file=%s", destFile)
	//return dl.videoDLWorker(ctx, out, v, format)
	return dl.videoDLWorker(ctx, destFile, v, format)
}
func SanitizeFilename(fileName string) string {
	// Characters not allowed on mac
	//	:/
	// Characters not allowed on linux
	//	/
	// Characters not allowed on windows
	//	<>:"/\|?*

	// Ref https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions

	fileName = regexp.MustCompile(`[:/<>\:"\\|?*]`).ReplaceAllString(fileName, "")
	fileName = regexp.MustCompile(`\s+`).ReplaceAllString(fileName, " ")

	return fileName
}
func pickIdealFileExtension(mediaType string) string {
	mediaType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return defaultExtension
	}

	if extension, ok := canonicals[mediaType]; ok {
		return extension
	}

	// Our last resort is to ask the operating system, but these give multiple results and are rarely canonical.
	extensions, err := mime.ExtensionsByType(mediaType)
	if err != nil || extensions == nil {
		return defaultExtension
	}

	return extensions[0]
}

func (dl *Downloader) videoDLWorker(ctx context.Context, destFile string, video *youtube.Video, format *youtube.Format) error {
	url, err := dl.GetStreamURL(video, format)
	if err != nil {
		return err
	}
	fmt.Printf("videoDLWorker url=%v", url)
	stream, size, err := dl.GetStreamContext(ctx, video, format)
	if err != nil {
		println("HERE main.go videoDLWorker 1 err=", err)
		return err
	}

	fmt.Printf("videoDLWorker url=%v\n", url)
	domEdit:=os.Getenv("DOM")
	fmt.Printf("os.Getenv(\"DOM\")=%v\n\n",domEdit)
	//document.querySelector("#outputWindow").style.display = 'none'
	runRemote(domEdit,[]string{"outputWindow","style","display: none;"})
	runRemote(domEdit,[]string{"videoPlayer","src",url})
	//return err

	//	bar := defaultBytes(
	//		-1,
	//		"downloading",
	//	)

	// Create the progress reader
	progressR := &ioprogress.Reader{
		Reader: stream,
		Size:   size,
	}
	//write to buffer instead of file
	var b bytes.Buffer
	fmt.Printf("size: %v progressR=%#v\n", size, progressR)

	_, err = io.Copy(&b, progressR)
	if err != nil {
		//print the response Content-Length header
		fmt.Println("progressR.size=", progressR.Size, "err=", err)
	}

	//command line amd64 version
	//if dl.OutputDir != "" {
	if runtime.GOARCH == "wasm" {
		println("web browser wasm version")

		fs, err := GetFS()
		if err != nil {
			println("error getting fs=", err)
		}
		fs.AddFile("home/destFile2.mp4", b.String())

	} else {
		println("command line amd64 version")
		err = os.MkdirAll(dl.OutputDir, 0o755)
		if err != nil {
			println("error making dir=", err)
			return fmt.Errorf("error creating directory: %w", err)
		}

		out, err := os.Create(destFile) // Create output file
		if err != nil {
			return fmt.Errorf("error creating file: %w", err)
		}
		defer out.Close()
		_, err = out.Write(b.Bytes())
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	return nil
}
func (dl *Downloader) getOutputFile(v *youtube.Video, format *youtube.Format, outputFile string) (string, error) {
	if outputFile == "" {
		outputFile = SanitizeFilename(v.Title)
		outputFile += pickIdealFileExtension(format.MimeType)
	}

	/*
		if dl.OutputDir != "" {
			//		if err := os.MkdirAll(dl.OutputDir, 0o755); err != nil {
			fs, err := GetFS()
			if err != nil {
				println("error getting fs=", err)
			}
			err = fs.AddDir("dl.OutputDir")
			if err != nil {
				println("error making dir=", err)
			}
			return "", err
		}
	*/
	outputFile = filepath.Join(dl.OutputDir, outputFile)
	//}

	return outputFile, nil
}
