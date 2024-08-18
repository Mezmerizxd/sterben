package youtube

import (
	"fmt"
	"os"
	"sterben/features/youtube"
	"sterben/pkg/pages"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// SetUrlPageModel represents the model for the "Set URL" page.
// It contains the configuration, text input model, metadata, and error states.
type SetUrlPageModel struct {
	Cfg             *pages.ModelConfig
	Input           textinput.Model
	InputError      string
	MetaData        *youtube.VideoMetaData
	MetaDataError   string
	MetaDataLoading bool
	Time            time.Time
}

// SetUrlPage initializes a new SetUrlPageModel with the provided configuration.
func SetUrlPage(cfg *pages.ModelConfig) *SetUrlPageModel {
	m := &SetUrlPageModel{
		Cfg:             cfg,
		MetaDataLoading: false,
		Time:            time.Now(),
	}

	// Initialize the text input with styles
	input := textinput.New()
	input.Placeholder = SetUrl.Name
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
func (p *SetUrlPageModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tick())
}

// Update handles incoming messages and updates the model state accordingly.
func (p *SetUrlPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	// Check if metadata is already loaded
	if p.MetaData != nil {
		return p.Cfg.Pages.SwitchToPreviousModel()
	}

	// Handle key messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if p.Input.Focused() {
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return p.Cfg.Pages.SwitchToPreviousModel()
			case tea.KeyEnter:
				p.InputError = ""
				if p.Input.Value() == "" {
					p.InputError = "Please enter a valid URL"
					return p, tea.Batch(cmds...)
				}

				// Start loading metadata in a goroutine
				p.MetaDataLoading = true
				go func() {
					metadata, err := youtube.GetVideoMetaData(p.Input.Value())
					if err != nil {
						p.MetaDataError = err.Error()
						p.MetaDataLoading = false
					} else {
						p.MetaData = metadata
						p.MetaDataLoading = false
					}
				}()

				return p, tea.Batch(cmds...)
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

// View renders the UI for the SetUrlPageModel.
func (p *SetUrlPageModel) View() string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	// Title
	title := lipgloss.NewStyle().Bold(true).Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(SetUrl.Name)

	// Input
	var input string
	if p.MetaDataLoading {
		input = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render("Loading...")
	} else {
		input = p.Input.View()
	}

	// Error handling
	err := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(p.InputError)
	if p.MetaDataError != "" {
		err = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(p.MetaDataError)
	}

	style := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(fmt.Sprintf("%s\n%s\n%s\n", title, input, err))
}

// Reset clears the input, metadata, and error states, resetting the page to its initial state.
func (p *SetUrlPageModel) Reset() {
	p.Input.Reset()
	p.MetaData = nil
	p.MetaDataError = ""
	p.MetaDataLoading = false
	p.InputError = ""
}
