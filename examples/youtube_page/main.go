package main

import (
	"fmt"
	"sterben/features/youtube"
	"sterben/pkg/config"
	"sterben/pkg/log"
	"sterben/pkg/pages"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	youtubeFeatureLog := log.New(log.Config{
		Feature:       "youtube_feature_test",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer youtubeFeatureLog.Close()
	youtubeFeature := NewYoutubePageModel(&pages.ModelConfig{
		Log:   youtubeFeatureLog,
		Pages: p,
	})

	// Start the Tea Program
	StartTeaProgram(p, youtubeFeature)
}

func StartTeaProgram(p *pages.Pages, feature *YoutubePageModel) {
	p.AddModel(YoutubePage, feature)
	p.SwitchModel(YoutubePage)

	tea := tea.NewProgram(p, tea.WithAltScreen(), tea.WithMouseAllMotion())

	_, err := tea.Run()
	if err != nil {
		p.Log.Error().Err(err).Msg("Failed to run tea program")
	}
}

var (
	YoutubePage pages.PageType = pages.PageType{
		ID:   "youtube_page",
		Name: "YoutubePage",
	}
)

type SetYoutubeUrl struct {
	input   textinput.Model
	err     string
	show    bool
	loading bool
}

type Options struct {
	list   []ListOption
	cursor ListOption
}

type YoutubePageModel struct {
	cfg             *pages.ModelConfig     // 8 bytes
	youtubeMetaData *youtube.VideoMetaData // 8 bytes
	options         Options                // place larger `Options` struct next
	setYoutubeUrl   SetYoutubeUrl          // place at the end
}

func NewYoutubePageModel(cfg *pages.ModelConfig) *YoutubePageModel {
	m := &YoutubePageModel{
		cfg: cfg,
		options: Options{
			list: []ListOption{setYoutubeUrl, other},
		},
	}

	m.options.cursor = m.options.list[0]

	ti := textinput.New()
	ti.Placeholder = "Youtube Url"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 100

	m.setYoutubeUrl.input = ti

	return m
}

func (p *YoutubePageModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p *YoutubePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if p.setYoutubeUrl.show {
		ti, cmd := p.setYoutubeUrl.input.Update(msg)
		p.setYoutubeUrl.input = ti
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		default:
			if !p.setYoutubeUrl.show {
				p.HandleOptionSelection(msg)
			} else {
				p.HandleInputSubmut(msg)
			}
			return p, tea.Batch(cmds...)
		}
	default:
		return p, tea.Batch(cmds...)
	}
}

func (p *YoutubePageModel) View() string {
	if p.setYoutubeUrl.show {
		return p.SetYoutubeUrlView()
	} else {
		return p.ListOptionsView()
	}
}

func (p *YoutubePageModel) SetYoutubeUrlView() string {
	var s string

	s += "Youtube Page / Set Url\n\n"

	if p.setYoutubeUrl.loading {
		s += "Loading...\n\n"
	}

	if p.setYoutubeUrl.err != "" {
		s += p.setYoutubeUrl.err + "\n\n"
	}

	if !p.setYoutubeUrl.loading {
		s += p.setYoutubeUrl.input.View()
	}
	return s
}

type ListOption string

const (
	setYoutubeUrl ListOption = "Set Youtube Url"
	other         ListOption = "Other"
)

func (p *YoutubePageModel) ListOptionsView() string {
	var s string

	s += "Youtube Page\n\n"

	if p.setYoutubeUrl.input.Value() != "" {
		s += fmt.Sprintf("Youtube Url: %s\n\n", p.setYoutubeUrl.input.Value())
	}

	for _, opt := range p.options.list {
		if opt == p.options.cursor {
			s += fmt.Sprintf("> %s\n", opt)
		} else {
			s += fmt.Sprintf("%s\n", opt)
		}
	}

	s += "\nPress Enter to select an option\n\n"

	if p.youtubeMetaData != nil {
		s += fmt.Sprintf("Found Video Data\nTitle: %s\nDescription: %s", p.youtubeMetaData.Title, p.youtubeMetaData.Description)
	}

	return s
}

func (p *YoutubePageModel) HandleOptionSelection(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyUp:
		for i, opt := range p.options.list {
			if opt == p.options.cursor && i > 0 {
				p.options.cursor = p.options.list[i-1]
				break
			}
		}
		return
	case tea.KeyDown:
		for i, opt := range p.options.list {
			if opt == p.options.cursor && i < len(p.options.list)-1 {
				p.options.cursor = p.options.list[i+1]
				break
			}
		}
		return
	case tea.KeyEnter:
		if p.options.cursor == setYoutubeUrl {
			p.setYoutubeUrl.show = true
		}
		return
	default:
		return
	}
}

func (p *YoutubePageModel) HandleInputSubmut(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyEnter:
		go func() {
			p.setYoutubeUrl.loading = true

			metaData, err := youtube.GetVideoMetaData(p.setYoutubeUrl.input.Value())
			if err != nil {
				p.setYoutubeUrl.err = err.Error()
			}

			p.youtubeMetaData = metaData

			p.setYoutubeUrl.loading = false

			p.setYoutubeUrl.show = false
		}()

		return
	}
}
