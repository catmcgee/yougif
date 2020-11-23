package service

import (
	"catmcgee/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func CreateVideoFrames(filePath, outputDirectory string, framesPerSecond string) error {
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		return err
	}

	command := exec.Command(config.Configuration.Processing.FfmpegBinaryPath, "-i", fmt.Sprintf("%s/%s", filePath, defaultVideoFileName), "-vf", fmt.Sprintf("fps=%s", framesPerSecond), fmt.Sprintf("%s/frame%%06d.jpg", outputDirectory))
	output, err := command.CombinedOutput()
	if err != nil {
		logrus.Println(string(output))
		return err
	}

	logrus.Println(string(output))
	return nil
}
