package tui

import (
	"fmt"
	"os"
	"sterben/pkg/pages"
	"sterben/tui/youtube"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// HomePageModel represents the model for the home page, containing configuration,
// options, alert message, and the current time.
type HomePageModel struct {
	Cfg     *pages.ModelConfig
	Options struct {
		List   []pages.PageType
		Cursor pages.PageType
	}
}

// HomePage initializes a new HomePageModel with the provided configuration.
func HomePage(cfg *pages.ModelConfig) *HomePageModel {
	m := &HomePageModel{
		Cfg: cfg,
	}

	m.Options.List = []pages.PageType{
		Youtube,
		ImageToIcon,
	}

	m.Options.Cursor = m.Options.List[0]

	return m
}

// Init is called when the program starts and returns the initial command.
func (p *HomePageModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model state accordingly.
func (p *HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyBackspace:
			return p.Cfg.Pages.SwitchToPreviousModel()
		default:
			m, cmd := p.handleOptions(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	default:
		return p, tea.Batch(cmds...)
	}
}

// View renders the UI for the HomePageModel.
func (p *HomePageModel) View() string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	// Title
	title := lipgloss.NewStyle().Bold(true).Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(Home.Name)

	// Options
	var options string
	for _, opt := range p.Options.List {
		if opt.ID == p.Options.Cursor.ID {
			options += "> " + opt.Name + "\n"
		} else {
			options += "  " + opt.Name + "\n"
		}
	}
	options = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(options)

	style := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(fmt.Sprintf(
		"%s\n\n%s",
		title,
		options,
	))
}

// handleOptions processes key messages for navigating and selecting options in the menu.
func (p *HomePageModel) handleOptions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		for i, opt := range p.Options.List {
			if opt.ID == p.Options.Cursor.ID && i > 0 {
				p.Options.Cursor = p.Options.List[i-1]
				break
			}
		}
		return p, nil

	case tea.KeyDown:
		for i, opt := range p.Options.List {
			if opt.ID == p.Options.Cursor.ID && i < len(p.Options.List)-1 {
				p.Options.Cursor = p.Options.List[i+1]
				break
			}
		}
		return p, nil

	case tea.KeyEnter:
		return p.handleEnter()

	default:
		return p, nil
	}
}

// handleEnter processes the action when the Enter key is pressed, initiating a download or navigating to another page.
func (p *HomePageModel) handleEnter() (tea.Model, tea.Cmd) {
	switch p.Options.Cursor {
	case Youtube:
		return p.Cfg.Pages.SwitchModel(youtube.Home)
	case ImageToIcon:
		return p.Cfg.Pages.SwitchModel(ImageToIcon)
	default:
		return p, nil
	}
}
