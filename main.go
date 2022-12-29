//go:build integration
// +build integration

package main

import (
	"context"
	"github.com/dougwatson/youtube"
)

func main() {
	ctx := context.Background()

	video, err := downloader.testDownloader.Client.GetVideoContext(ctx, "BaW_jenozKc")

	downloader.testDownloader.DownloadComposite(ctx, "", video, "hd1080", "mp4")
}
