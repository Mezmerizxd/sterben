package main

import (
	"sterben/pkg/config"
	"sterben/pkg/log"
	"sterben/pkg/pages"
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

	// Initialize Pages
	p := pages.Initialize(pages.Config{
		Log: log.New(log.Config{
			Feature:       "pages",
			ConsoleOutput: false,
			FileOutput:    true,
		}),
	})
	defer p.Close()
}
