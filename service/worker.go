package service

import (
	"catmcgee/config"
	"catmcgee/model"
	"catmcgee/repository"
	ytHelper "catmcgee/youtube_helper"
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"html"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	framesPerSecond = 1
)

func GetYouTubeVideos(region string, totalNumberOfVideos int) {
	ctx := context.Background()

	clientSecretFile, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		logrus.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(clientSecretFile, youtube.YoutubeReadonlyScope)
	if err != nil {
		logrus.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := ytHelper.GetClient(ctx, config)
	service, err := youtube.New(client)

	ytHelper.HandleError(err, "Error creating YouTube client")
	categoryIds, err := ytHelper.GetCategories(service, region)
	if err != nil {
		logrus.Printf("unable to get categories from youttube: %s", err.Error())
		return
	}

	videosFound := 0
	for _, categoryId := range categoryIds {
		if videosFound >= totalNumberOfVideos {
			continue
		}
		videosFound = videosFound + ytHelper.VideoList(service, []string{"id", "snippet", "contentDetails"}, region, "", categoryId, totalNumberOfVideos, 0)
	}
}

func StartProcessingVideos(numberOfWorker int) {
	jobChannel := make(chan *model.Video)

	for i := 0; i < numberOfWorker; i++ {
		go startWorker(jobChannel)
	}

	getAndProcessVideos(jobChannel)
}

func getAndProcessVideos(jobChannel chan *model.Video) {
	videos, err := repository.SelectCaptionAndNotProcessedVideos(context.Background())
	if err != nil {
		logrus.Println(err)
		return
	}

	for _, video := range videos {
		jobChannel <- video
	}

	time.Sleep(30 * time.Second)
	getAndProcessVideos(jobChannel)
}

func startWorker(jobsChannel <-chan *model.Video) {
	for job := range jobsChannel {
		if err := processVideoWithId(job); err != nil {
			logrus.Println(err)
		}
	}
}

func processVideoWithId(video *model.Video) error {
	logrus.Printf("Start processing video: '%s'", video.Id)

	tx, err := repository.Begin()
	if err != nil {
		return err
	}

	captionTrack, err := getVideoCaption(video.VideoId)
	if err != nil {
		video.HasCaption = false
		repository.UpdateVideo(context.Background(), tx, video)
		tx.Commit()
		return err
	}

	transcript, err := getCaptions(captionTrack.BaseUrl)
	if err != nil {
		tx.Rollback()
		return err
	}

	textPartsOnly := make([]string, 0, len(transcript.Text))

	captions := make([]*model.Caption, 0, len(transcript.Text))
	for _, text := range transcript.Text {
		decodedText := html.UnescapeString(text.Text)

		startInMs := int(text.Start * 1000)
		endInMs := startInMs + int(text.Duration*1000)
		captions = append(captions, model.NewCaption(video.Id, decodedText, startInMs, endInMs))

		textPartsOnly = append(textPartsOnly, decodedText)
	}

	video.Caption = strings.Join(textPartsOnly, " ")
	if err := repository.UpdateVideo(context.Background(), tx, video); err != nil {
		tx.Rollback()
		return err
	}

	for i, caption := range captions {
		var previousCaptionId *string
		if i > 0 {
			previousCaptionId = &captions[i-1].Id
		}

		var nextCaptionId *string
		var nextStart int
		if i < len(captions)-1 {
			nextCaption := captions[i+1]
			nextCaptionId = &nextCaption.Id
			nextStart = nextCaption.Start
		}

		caption.PreviousCaption = previousCaptionId
		caption.NextCaption = nextCaptionId

		if (caption.Start == caption.End) && nextStart >= 0 {
			caption.End = nextStart
		}
		if err := repository.InsertCaption(context.Background(), tx, caption); err != nil {
			logrus.Println(err)
		}
	}

	videoDirectory, err := DownloadVideoWithId(video.VideoId)
	if err != nil {
		tx.Rollback()
		return err
	}

	defer cleanUpWorkingDirectory(videoDirectory)

	framesDirectory := fmt.Sprintf("%s/frames", videoDirectory)
	if err := CreateVideoFrames(videoDirectory, framesDirectory, config.Configuration.Processing.FramesPerSecond); err != nil {
		tx.Rollback()
		return err
	}

	fileInfos, err := ioutil.ReadDir(framesDirectory)
	if err != nil {
		tx.Rollback()
		return err
	}

	previousFrameId := ""
	currentTime := 0

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", framesDirectory, fileInfo.Name())
		frameId, err := storeFrameIntoDatabase(tx, filePath, fileInfo.Name(), video.Id, previousFrameId, currentTime)
		if err != nil {
			logrus.Println(err)
		}
		previousFrameId = frameId
		currentTime += 1000 / framesPerSecond
	}

	video.Processed = true
	if err := repository.UpdateVideo(context.Background(), tx, video); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func storeFrameIntoDatabase(tx *sql.Tx, filePath, fileName, videoId, previousFrameId string, currentTime int) (string, error) {
	imageFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	imageContent, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return "", err
	}

	frame := model.NewFrame(videoId, fileName, currentTime, imageContent, previousFrameId)
	if err := repository.InsertFrame(context.Background(), tx, frame); err != nil {
		return "", err
	}

	return frame.Id, nil
}

func cleanUpWorkingDirectory(directory string) {
	if err := os.RemoveAll(directory); err != nil {
		logrus.Println(err)
	}
}
