package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/knadh/go-get-youtube/youtube"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"os/exec"
	"time"
	"yttomp3/downloader"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("could not initialize configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("could not load env variables: %s", err.Error())
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	logrus.Println("Bot started successfully")

	b.Handle("/start", func (m *tb.Message){
		b.Send(m.Sender, "Hello, I can download your music from YouTube. Just send me a link to any music video on YouTube")
	})

	if err != nil {
		logrus.Fatalf("unable to launch bot: %s",err.Error())
		return
	}

	b.Handle(tb.OnText, func(m *tb.Message) {

		d, err:= downloader.NewDownloader(viper.GetInt64("max_video_duration"), viper.GetInt64("max_download_time"))
		if err != nil {
			return
		}
		u := m.Text
		if !d.IsValidUrl(u){
			b.Send(m.Sender, "message is not valid url. Send link to a video")
			return
		}
		logrus.Println("Starting the download...")
		filename, err := d.Download(u)
		if err != nil{
			b.Send(m.Sender, fmt.Sprintf("was unable to download audio.\nError Description: %s", err.Error()))
		}
		logrus.Println("Finished")
		file := tb.FromDisk(filename)
		fmt.Println(filename)

		audio := &tb.Audio{File: file, FileName: filename, Title: filename}
		audio.Send(b, m.Sender, &tb.SendOptions{ReplyTo: m})

		cmd := exec.Command("rm", filename)
		cmd.CombinedOutput()
	})

	b.Start()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}