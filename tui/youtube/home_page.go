package youtube

import (
	"fmt"
	"os"
	"sterben/features/youtube"
	"sterben/pkg/pages"
	"time"

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
	Alert string
	Time  time.Time
}

// tickMsg is a custom message used to update the time every second.
type tickMsg time.Time

// downloadMsg is a custom message used to signal the result of the download process.
type downloadMsg struct {
	success bool
	err     error
}

// clearAlertMsg is a custom message used to clear the alert after a certain duration.
type clearAlertMsg struct{}

// HomePage initializes a new HomePageModel with the provided configuration.
func HomePage(cfg *pages.ModelConfig) *HomePageModel {
	m := &HomePageModel{
		Cfg:  cfg,
		Time: time.Now(),
	}

	m.Options.List = []pages.PageType{
		SetUrl,
		Download,
	}

	m.Options.Cursor = m.Options.List[0]

	return m
}

// tick returns a command that sends a tickMsg every second.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init is called when the program starts and returns the initial command.
func (p *HomePageModel) Init() tea.Cmd {
	return tick()
}

// Update handles incoming messages and updates the model state accordingly.
func (p *HomePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		p.Time = time.Time(msg)
		return p, tick()

	case downloadMsg:
		if msg.err != nil {
			p.Alert = msg.err.Error()
		} else if msg.success {
			p.Alert = "Downloaded!"
		}
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(3 * time.Second)
			return clearAlertMsg{}
		})

	case clearAlertMsg:
		p.Alert = ""

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

	return p, tea.Batch(cmds...)
}

// View renders the UI for the HomePageModel.
func (p *HomePageModel) View() string {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	// Title
	title := lipgloss.NewStyle().Bold(true).Padding(1).Foreground(lipgloss.Color("#ff1f1f")).Render(Home.Name)

	// Alert
	var alert string
	if p.Alert != "" {
		alert = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff1f1f")).Render(p.Alert)
	}

	// Options
	var options string
	for _, opt := range p.Options.List {
		if opt.ID == p.Options.Cursor.ID {
			options += "> "
		} else {
			options += "  "
		}

		switch opt {
		case SetUrl:
			setUrlPageModel := p.Cfg.Pages.Models[SetUrl].(*SetUrlPageModel)
			if setUrlPageModel.MetaData != nil {
				options += lipgloss.NewStyle().Foreground(lipgloss.Color("#ff1f1f")).Render(opt.Name, "(Reset)") + "\n"
			} else {
				options += opt.Name + "\n"
			}
		default:
			options += opt.Name + "\n"
		}
	}
	options = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(options)

	style := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render(fmt.Sprintf(
		"%s\n%s\n\n%s\n%s\n",
		title,
		fmt.Sprintf("Current time: %v", p.Time.Format("15:04:05")),
		alert,
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
	case SetUrl:
		setUrlPageModel := p.Cfg.Pages.Models[SetUrl].(*SetUrlPageModel)
		setUrlPageModel.Reset()
		return p.Cfg.Pages.SwitchModel(SetUrl)

	case Download:
		setUrlPageModel := p.Cfg.Pages.Models[SetUrl].(*SetUrlPageModel)

		if setUrlPageModel.MetaData == nil {
			p.Alert = "No metadata available"
			return p, tea.Batch(func() tea.Msg {
				time.Sleep(3 * time.Second)
				return clearAlertMsg{}
			})
		}

		p.Alert = "Downloading..."
		cmds := []tea.Cmd{
			func() tea.Msg {
				err := youtube.DownloadYoutubeVideo(setUrlPageModel.Input.Value(), "downloads")
				if err != nil {
					return downloadMsg{err: err}
				}
				return downloadMsg{success: true}
			},
		}
		return p, tea.Batch(cmds...)

	default:
		return p, nil
	}
}
