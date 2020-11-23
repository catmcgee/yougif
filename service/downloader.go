package service

import (
	"context"
	"fmt"
	"github.com/rylio/ytdl"
	"net/url"
	"os"
	"strings"
)

const (
	defaultVideoFileName = "video.mp4"

	defaultResolution = "720p"
	defaultExtension  = "mp4"
)

var youTubeClient = ytdl.DefaultClient

func DownloadVideoWithId(id string) (string, error) {
	return DownloadVideoWithUrl(fmt.Sprintf("https://www.youtube.com/watch?v=%s", id))
}

func DownloadVideoWithUrl(videoUrl string) (string, error) {
	videoInfo, err := youTubeClient.GetVideoInfo(context.Background(), videoUrl)
	if err != nil {
		return "", err
	}

	uri, err := url.ParseRequestURI(videoUrl)
	if err != nil {
		return "", err
	}

	videoId := extractVideoID(uri)

	format := getFormatByResolutionAndExtension(videoInfo.Formats, defaultResolution, defaultExtension)
	if format == nil {
		return "", fmt.Errorf("none of the following resolutions found [1080p, 720p, 360p] for video '%s'", videoId)
	}

	videoDirectory := fmt.Sprintf("downloads/%s", videoId)
	if err := os.MkdirAll(videoDirectory, os.ModePerm); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", videoDirectory, defaultVideoFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if err := youTubeClient.Download(context.Background(), videoInfo, format, file); err != nil {
		// TODO: cleanup in case of error
		return "", err
	}

	return videoDirectory, nil
}

func extractVideoID(u *url.URL) string {
	switch u.Host {
	case "www.youtube.com", "youtube.com", "m.youtube.com":
		if u.Path == "/watch" {
			return u.Query().Get("v")
		}
		if strings.HasPrefix(u.Path, "/embed/") {
			return u.Path[7:]
		}
	case "youtu.be":
		if len(u.Path) > 1 {
			return u.Path[1:]
		}
	}
	return ""
}

func getFormatByResolutionAndExtension(formatList ytdl.FormatList, resolution, extension string) *ytdl.Format {
	for _, format := range formatList {
		if format.Resolution == resolution && format.Extension == extension {
			return format
		}
	}
	//
	//if format := getFormatByResolutionAndExtension(formatList, "720p", extension); format != nil {
	//	return format
	//}

	if format := getFormatByResolutionAndExtension(formatList, "360p", extension); format != nil {
		return format
	}

	return nil
}
