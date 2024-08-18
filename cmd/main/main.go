package main

import (
	"fmt"
	"sterben/features/youtube"
	"sterben/pkg/config"
	"sterben/pkg/log"
	"sterben/tui"
)

func main() {
	// Initialize Log
	l := log.New(log.Config{
		Feature:       "main",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer l.Close()

	// Initialize Config
	config.Initialize()
	cfg, err := config.GetConfig()
	if err != nil {
		l.Error().Err(err).Msg("Failed to get config")
		panic(err)
	}
	if cfg != nil {
		l.Info().Msg("Config loaded successfully")
	} else {
		l.Info().Msg("Config not loaded")
	}

	if !youtube.CheckIfYtdlpInstalled() {
		fmt.Println("yt-dlp is not installed, installing...	")
		// Attempt to install yt-dlp
		err := youtube.DownloadYtdlp()
		if err != nil {
			l.Error().Err(err).Msg("Failed to download yt-dlp")
			panic(err)
		}
	}

	y, err := tui.Initialize()
	if err != nil {
		l.Error().Err(err).Msg("Failed to initialize TUI")
		panic(err)
	}

	err = y.Start()
	if err != nil {
		l.Error().Err(err).Msg("Failed to start TUI")
		panic(err)
	}
}
