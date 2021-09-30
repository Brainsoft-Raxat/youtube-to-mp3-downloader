package configs

type Configs struct {
	MaxVideoDuration int //minutes
	MaxDownloadTime int //seconds
	BotUsername string
}

func NewConfigs(maxVideoDuration int, maxDownloadTime int, botUsername string) *Configs {
	return &Configs{MaxVideoDuration: maxVideoDuration, MaxDownloadTime: maxDownloadTime, BotUsername: botUsername}
}


