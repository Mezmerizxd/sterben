package start

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"sterben/pkg/pages"
)

const (
	StartPage pages.PageType = "start"
)

type StartPageModel struct {
	cfg     *pages.ModelConfig
	options []StartPageOption
	cursor  StartPageOption
	counter int
}

type StartPageOption struct {
	id       string
	label    string
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
			id:     "counter",
			label:  "Enter Counter",
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
				if opt.id == p.cursor.id && i > 0 {
					p.cursor = p.options[i-1]
					break
				}
			}
			return p, nil
		case tea.KeyDown:
			for i, opt := range p.options {
				if opt.id == p.cursor.id && i < len(p.options)-1 {
					p.cursor = p.options[i+1]
					break
				}
			}
			return p, nil
		case tea.KeyEnter:
			if p.cursor.isPage {
				p.cfg.Log.Info().Msgf("Switching to %s Page", p.cursor.label)
				return p.cfg.Pages.SwitchModel(pages.PageType(p.cursor.id))
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
		if opt.id == p.cursor.id {
			s += fmt.Sprintf("%s\n", spm_optionStyle(">", opt.label))
		} else {
			s += fmt.Sprintf("%s\n", spm_textStyle(" ", opt.label))
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
