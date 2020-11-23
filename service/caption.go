package service

import (
	"catmcgee/model"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	videoInfoBaseUrl = "https://youtube.com/get_video_info?video_id=%s"
)

var (
	captionTrackRegex = regexp.MustCompile(`{"captionTracks":.*isTranslatable":(true|false)}]`)
)

func getVideoCaption(videoId string) (*model.CaptionTrack, error) {
	response, err := http.Get(fmt.Sprintf(videoInfoBaseUrl, videoId))
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return nil, fmt.Errorf("video info response status code is '%d'", response.StatusCode)
	}
	defer response.Body.Close()

	encodedContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	decodedContent, err := url.QueryUnescape(string(encodedContent))
	if err != nil {
		return nil, err
	}

	foundString := captionTrackRegex.FindString(decodedContent)
	if len(foundString) == 0 {
		return nil, fmt.Errorf("no caption found for video '%s'", videoId)
	}

	final := fmt.Sprintf("%s}", foundString)

	var captions model.Captions
	if err := json.Unmarshal([]byte(final), &captions); err != nil {
		return nil, err
	}

	for _, captionTrack := range captions.CaptionTracks {
		if strings.Contains(captionTrack.VssId, "en") {
			return &captionTrack, nil
		}
	}

	return nil, fmt.Errorf("no caption found for video '%s'", videoId)
}

func getCaptions(baseUrl string) (*model.Transcript, error) {
	response, err := http.Get(baseUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	xmlContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var transcript model.Transcript
	if err := xml.Unmarshal(xmlContent, &transcript); err != nil {
		return nil, err
	}

	return &transcript, nil
}

func getCaptionsWithSearchString(captions []*model.Caption, searchString string) []*model.Caption {
	result := make([]*model.Caption, 0)

	searchString = strings.ToLower(searchString)

	textOnly := make([]string, 0, len(captions))
	for _, caption := range captions {
		textOnly = append(textOnly, strings.ToLower(caption.Text))
	}

	text := strings.Join(textOnly, " ")
	searchStartIndex := strings.Index(text, searchString)
	if searchStartIndex == -1 {
		return nil
	}

	searchEndIndex := searchStartIndex + len(searchString)

	textLength := 0
	for _, caption := range captions {
		captionStartIndex := textLength
		captionEndIndex := captionStartIndex + len(caption.Text)

		textLength = captionEndIndex + 1

		if searchStartIndex >= captionStartIndex && searchStartIndex <= captionEndIndex {
			result = append(result, caption)
			continue
		}

		if searchEndIndex >= captionStartIndex && searchEndIndex <= captionEndIndex {
			result = append(result, caption)
		}
	}

	return result
}
