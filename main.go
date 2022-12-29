package main

import (
	"context"

	"github.com/dougwatson/youtube/v2"
)

func main() {
	ctx := context.Background()

	downloader := NewDownloader()
	youtube.RegisterDownloader(downloader)

	video, err := downloader.testDownloader.Client.GetVideoContext(ctx, "BaW_jenozKc")

	downloader.testDownloader.DownloadComposite(ctx, "", video, "hd1080", "mp4")
}
