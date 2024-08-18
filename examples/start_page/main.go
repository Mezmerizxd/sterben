package main

import (
	"fmt"
	"sterben/pkg/config"
	"sterben/pkg/log"
	"sterben/pkg/pages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Initialize Log
	l := log.New(log.Config{
		Feature:       "main_test",
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
			Feature:       "pages_test",
			ConsoleOutput: false,
			FileOutput:    true,
		}),
	})
	defer p.Close()

	// Initalize Features
	startFeatureLog := log.New(log.Config{
		Feature:       "start_feature_test",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer startFeatureLog.Close()
	startFeature := NewStartPageModel(&pages.ModelConfig{
		Log:   startFeatureLog,
		Pages: p,
	})

	// Start the Tea Program
	StartTeaProgram(p, startFeature)
}

func StartTeaProgram(p *pages.Pages, startFeature *StartPageModel) {
	p.AddModel(StartPage, startFeature)
	p.SwitchModel(StartPage)

	tea := tea.NewProgram(p, tea.WithAltScreen(), tea.WithMouseAllMotion())

	_, err := tea.Run()
	if err != nil {
		p.Log.Error().Err(err).Msg("Failed to run tea program")
	}
}

var StartPage pages.PageType = pages.PageType{
	ID:   "start",
	Name: "Start",
}

type StartPageModel struct {
	cfg     *pages.ModelConfig
	options []StartPageOption
	cursor  StartPageOption
	counter int
}

type StartPageOption struct {
	page     pages.PageType
	isPage   bool
	function func()
}

func NewStartPageModel(cfg *pages.ModelConfig) *StartPageModel {
	m := &StartPageModel{
		cfg:     cfg,
		counter: 0,
	}

	m.options = []StartPageOption{
		{
			page:   StartPage,
			isPage: false,
			function: func() {
				m.counter++
			},
		},
	}

	m.cursor = m.options[0]

	return m
}

func (p *StartPageModel) Init() tea.Cmd {
	return nil
}

func (p *StartPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		case tea.KeyUp:
			for i, opt := range p.options {
				if opt.page.ID == p.cursor.page.ID && i > 0 {
					p.cursor = p.options[i-1]
					break
				}
			}
			return p, nil
		case tea.KeyDown:
			for i, opt := range p.options {
				if opt.page.ID == p.cursor.page.ID && i < len(p.options)-1 {
					p.cursor = p.options[i+1]
					break
				}
			}
			return p, nil
		case tea.KeyEnter:
			if p.cursor.isPage {
				p.cfg.Log.Info().Msgf("Switching to %s Page", p.cursor.page.Name)
				return p.cfg.Pages.SwitchModel(pages.PageType(p.cursor.page))
			} else {
				p.cursor.function()
				return p, nil
			}
		default:
			return p, nil
		}
	default:
		return p, nil
	}
}

func (p *StartPageModel) View() string {
	s := fmt.Sprintf("%s\n", spm_titleStyle("Sterben"))

	// Counter
	s += fmt.Sprintf("%s\n", spm_textStyle(fmt.Sprintf("Counter: %d", p.counter)))

	for _, opt := range p.options {
		if opt.page.ID == p.cursor.page.ID {
			s += fmt.Sprintf("%s\n", spm_optionStyle(">", opt.page.Name))
		} else {
			s += fmt.Sprintf("%s\n", spm_textStyle(" ", opt.page.Name))
		}
	}

	// If no Options are available, show a message
	if len(p.options) == 0 {
		s += spm_textStyle("No options available")
	}

	return s
}

var (
	spm_titleStyle  = lipgloss.NewStyle().Foreground(pages.TextPrimaryStyle).Bold(true).Padding(1).Render
	spm_textStyle   = lipgloss.NewStyle().Foreground(pages.TextSecondaryStyle).Padding(1).Render
	spm_optionStyle = lipgloss.NewStyle().Foreground(pages.TextPrimaryStyle).Bold(true).Padding(1).Render
)
