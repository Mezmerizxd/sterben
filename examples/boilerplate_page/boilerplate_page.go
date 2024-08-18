package boilerplate_page

import (
	"sterben/pkg/pages"

	tea "github.com/charmbracelet/bubbletea"
)

type BoilerplatePageModel struct {
	cfg *pages.ModelConfig
}

func BoilerplatePage(cfg *pages.ModelConfig) *BoilerplatePageModel {
	m := &BoilerplatePageModel{
		cfg: cfg,
	}

	return m
}

func (p *BoilerplatePageModel) Init() tea.Cmd {
	return nil
}

func (p *BoilerplatePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
		default:
			return p, tea.Batch(cmds...)
		}
	default:
		return p, tea.Batch(cmds...)
	}
}

func (p *BoilerplatePageModel) View() string {
	return "Boilerplate Page"
}
