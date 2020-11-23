package youtube_helper

import (
	"catmcgee/model"
	"catmcgee/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func GetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		logrus.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		logrus.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		logrus.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logrus.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		logrus.Fatalf(message+": %v", err.Error())
	}
}

func GetCategories(service *youtube.Service, region string) ([]string, error) {
	videoCategoriesListCall := service.VideoCategories.List([]string{"id"})
	videoCategoriesListCall.RegionCode(region)
	response, err := videoCategoriesListCall.Do()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(response.Items))
	for _, item := range response.Items {
		result = append(result, item.Id)
	}

	return result, nil
}

func VideoList(service *youtube.Service, parts []string, region, nextPageToken, categoryId string, totalNumberOfVideos, numberOfVideos int) int {
	if numberOfVideos >= totalNumberOfVideos {
		return numberOfVideos
	}

	videoListCall := service.Videos.List(parts)
	videoListCall.Chart("mostPopular")
	videoListCall.VideoCategoryId(categoryId)
	videoListCall.RegionCode(region)
	videoListCall.PageToken(nextPageToken)
	videoListCall.MaxResults(25)
	response, err := videoListCall.Do()
	if err != nil {
		logrus.Println(err)
		return numberOfVideos
	}

	for _, item := range response.Items {
		logrus.WithFields(logrus.Fields{
			"videoId":    item.Id,
			"title":      item.Snippet.Title,
			"hasCaption": item.ContentDetails.Caption,
		}).Println("found video")
		hasCaption, err := strconv.ParseBool(item.ContentDetails.Caption)
		if err != nil {
			logrus.Println(err)
			hasCaption = false
		}

		video := model.NewVideo(item.Id, "", hasCaption)
		if err := repository.InsertVideo(context.Background(), video); err != nil {
			logrus.Println(err)
		}
	}

	if int64(numberOfVideos) < response.PageInfo.TotalResults {
		return VideoList(service, parts, region, response.NextPageToken, categoryId, totalNumberOfVideos, numberOfVideos+len(response.Items))
	}

	return numberOfVideos
}
