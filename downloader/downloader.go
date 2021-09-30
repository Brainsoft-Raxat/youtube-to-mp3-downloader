package downloader

import (
	"context"
	"errors"
	"fmt"
	yt "github.com/kkdai/youtube/v2"
	_ "github.com/spf13/viper"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	expression = "^((?:https?:)?\\/\\/)?((?:www|m)\\.)?((?:youtube\\.com|youtu.be))(\\/(?:[\\w\\-]+\\?v=|embed\\/|v\\/)?)([\\w\\-]+)(\\S+)?$"
)

type Downloader struct {
	MaxVideoDuration time.Duration
	MaxDownloadTime time.Duration
	r                *regexp.Regexp
}

func NewDownloader(maxVideoDuration, maxDownloadTime int64) (*Downloader, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	return &Downloader{
		MaxVideoDuration: time.Minute * time.Duration(maxVideoDuration),
		MaxDownloadTime: time.Second * time.Duration(maxDownloadTime),
		r:                r,
	}, nil
}

func (d *Downloader) Download(u string) (string, error) {
	client := yt.Client{}

	video, err := client.GetVideo(u)
	if err != nil {
		return "", err
	}

	if video.Duration > d.MaxVideoDuration {
		return "", errors.New(fmt.Sprintf("video duration exceeds maximum limit of %d minutes", d.MaxVideoDuration))
	}

	filename := video.Title
	filename = strings.Replace(filename, " ", "_", -1)
	_ , cancel := context.WithTimeout(context.Background(), time.Second*d.MaxDownloadTime)
	defer cancel()
	// youtube-dl -x --extract-audio --audio-format mp3 <video URL>
	//cmd := exec.CommandContext(ctx, "youtube-dl", "-x", "--extract-audio", "--audio-format", "mp3", u, "-o", filename+".%(ext)s")
	cmd := exec.Command("youtube-dl", "-o", filename+".%(ext)s", "-x", "--audio-format", "mp3", u)
	data, err := cmd.CombinedOutput()

	if err != nil {
		os.Remove(filename)
		return "", errors.New(fmt.Sprintf("%s error from CombinedOutput, data: %s", err.Error(), data))
	}

	if strings.Contains(string(data), "ERROR") {
		os.Remove(filename)
		return "", errors.New(fmt.Sprintf("error downloading video with youtube-dl, output: %s", string(data)))
	}

	return filename + ".mp3", nil
}

func (d *Downloader) IsValidUrl(u string) bool {
	return d.r.MatchString(u)
}

