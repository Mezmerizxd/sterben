package tui

import (
	"fmt"
	"os"
	"sterben/features/image_convert"
	"sterben/pkg/pages"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// ImageToIconPageModel represents the model for the "Set URL" page.
// It contains the configuration, text input model, metadata, and error states.
type ImageToIconPageModel struct {
	Cfg                *pages.ModelConfig
	Input              textinput.Model
	InputError         string
	ImageToIconError   string
	ImageToIconLoading bool
	Time               time.Time
}

// ImageToIconPage initializes a new ImageToIconPageModel with the provided configuration.
func ImageToIconPage(cfg *pages.ModelConfig) *ImageToIconPageModel {
	m := &ImageToIconPageModel{
		Cfg:                cfg,
		ImageToIconLoading: false,
		Time:               time.Now(),
	}

	// Initialize the text input with styles
	input := textinput.New()
	input.Placeholder = "Enter Image Path"
	input.Focus()

	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff1f1f"))
	input.Cursor.Style = lipgloss.NewStyle().Background(lipgloss.Color("#ff1f1f"))
	input.Cursor.TextStyle = redStyle
	input.TextStyle = redStyle
	input.CompletionStyle = redStyle
	input.PlaceholderStyle = redStyle
	input.PromptStyle = redStyle

	m.Input = input
	return m
}

// Init initializes the model, setting up the blinking cursor for text input.
func (p *ImageToIconPageModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tick())
}

// Update handles incoming messages and updates the model state accordingly.
func (p *ImageToIconPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update the time
	switch msg := msg.(type) {
	case tickMsg:
		p.Time = time.Time(msg)
		cmds = append(cmds, tick())
	}

	// Update text input
	ti, cmd := p.Input.Update(msg)
	p.Input = ti
	cmds = append(cmds, cmd)

	// Handle key messages
	switch msg := msg.(type) {
	case switchPrevPageMsg:
		return p.Cfg.Pages.SwitchToPreviousModel()
	case tea.KeyMsg:
		if p.Input.Focused() {
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return p.Cfg.Pages.SwitchToPreviousModel()
			case tea.KeyEnter:
				p.InputError = ""
				if p.Input.Value() == "" {
					p.InputError = "Please enter a valid Path"
					return p, tea.Batch(cmds...)
				}

				// Start loading metadata in a goroutine
				p.ImageToIconLoading = true
				go func() {
					err := image_convert.ConvertImageToMultipleIcons(p.Input.Value(), image_convert.Sizes)
					if err != nil {
						p.ImageToIconError = err.Error()
						p.ImageToIconLoading = false
					}
				}()

				return p, tea.Batch(func() tea.Msg {
					time.Sleep(2 * time.Second)
					p.ImageToIconLoading = false
					if p.ImageToIconError != "" {
						return tea.Batch(cmds...)
					}
					return switchPrevPageMsg{}
				})
			}
		} else {
			// Handle non-focused input key messages
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc, tea.KeyBackspace:
				return p.Cfg.Pages.SwitchToPreviousModel()
			}
		}

	default:
		return p, tea.Batch(cmds...)
	}

	return p, tea.Batch(cmds...)
}

// View renders the UI for the ImageToIconPageModel.
func (p *ImageToIconPageModel) View() string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	// Title
	title := lipgloss.NewStyle().Bold(true).Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(ImageToIcon.Name)

	// Input
	var input string
	if p.ImageToIconLoading {
		input = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render("Loading...")
	} else {
		input = p.Input.View()
	}

	// Error handling
	err := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(p.InputError)
	if p.ImageToIconError != "" {
		err = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(p.ImageToIconError)
	}

	style := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(fmt.Sprintf("%s\n%s\n%s", title, input, err))
}

// Reset clears the input, metadata, and error states, resetting the page to its initial state.
func (p *ImageToIconPageModel) Reset() {
	p.Input.Reset()
	p.ImageToIconError = ""
	p.ImageToIconLoading = false
	p.InputError = ""
}

// tickMsg is a custom message used to update the time every second.
type tickMsg time.Time

// tick returns a command that sends a tickMsg every second.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type switchPrevPageMsg struct{}
