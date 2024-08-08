package main

import (
	"sterben/pkg/config"
	"sterben/pkg/features/start"
	"sterben/pkg/log"
	"sterben/pkg/pages"

	tea "github.com/charmbracelet/bubbletea"
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

	// Initalize Features
	startFeatureLog := log.New(log.Config{
		Feature:       "start_feature",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer startFeatureLog.Close()
	startFeature := start.NewStartPageModel(&pages.ModelConfig{
		Log:   startFeatureLog,
		Pages: p,
	})

	// Start the Tea Program
	StartTeaProgram(p, startFeature)
}

func StartTeaProgram(p *pages.Pages, startFeature *start.StartPageModel) {
	p.AddModel(start.StartPage, startFeature)
	p.SwitchModel(start.StartPage)

	tea := tea.NewProgram(p, tea.WithAltScreen(), tea.WithMouseAllMotion())

	_, err := tea.Run()
	if err != nil {
		p.Log.Error().Err(err).Msg("Failed to run tea program")
	}
}
