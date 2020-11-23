package config

import (
	"encoding/json"
	"os"
)

var (
	Configuration *Config
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Db         DbConfig         `json:"database"`
	YouTube    YouTubeConfig    `json:"youTube"`
	Processing ProcessingConfig `json:"processing"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type DbConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	DatabaseName string `json:"databaseName"`
}

type YouTubeConfig struct {
	Region                         string `json:"region"`
	TotalNumbersOfVideosToDownload int    `json:"totalNumbersOfVideosToDownload"`
}

type ProcessingConfig struct {
	NumberOfWorkers  int    `json:"numberOfWorkers"`
	FfmpegBinaryPath string `json:"ffmpegBinaryPath"`
	FramesPerSecond  string `json:"framesPerSecond"`
}

func Load() error {
	file, err := os.Open("./config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return err
	}

	Configuration = &config
	return nil
}
