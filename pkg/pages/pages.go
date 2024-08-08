package pages

import (
	"errors"
	"sterben/pkg/log"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Singleton instance of Pages.
	instance *Pages
)

// GetInstance returns the singleton instance of Pages.
func GetInstance() *Pages {
	return instance
}

// Config holds the configuration required to initialize the Pages.
type Config struct {
	Log log.Log
}

// ModelConfig holds the configuration passed to each model.
type ModelConfig struct {
	Log   log.Log
	Pages *Pages
}

// PageType represents the type of a page, used as a key in the Models map.
type PageType string

// Pages manages the application's pages, including navigation and model handling.
type Pages struct {
	Log        log.Log
	Mutex      sync.Mutex
	Models     map[PageType]tea.Model
	Navigation []tea.Model
}

// Initialize creates a new Pages instance and stores it as a singleton.
func Initialize(c Config) *Pages {
	p := &Pages{
		Log:    c.Log,
		Models: make(map[PageType]tea.Model),
	}
	instance = p
	return p
}

// Init initializes the current model in the navigation stack.
func (p *Pages) Init() tea.Cmd {
	m, err := p.CurrentModel()
	if err != nil {
		p.Log.Error().Err(*err).Msg("Failed to get current model")
		return nil
	}

	return tea.Batch(
		tea.ClearScreen,
		m.Init(),
	)
}

// Close closes the logging instance.
func (p *Pages) Close() {
	p.Log.Close()
}

// Update handles updates to the current model based on messages received.
func (p *Pages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, err := p.CurrentModel()
	if err != nil {
		p.Log.Error().Err(*err).Msg("Failed to get current model")
		return nil, tea.Quit
	}

	return m.Update(msg)
}

// View renders the view of the current model.
func (p *Pages) View() string {
	m, err := p.CurrentModel()
	if err != nil {
		p.Log.Error().Err(*err).Msg("Failed to get current model")
		return ""
	}

	return m.View()
}

// AddModel adds a new model to the Pages instance and logs the action.
func (p *Pages) AddModel(t PageType, m tea.Model) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.Models[t] = m
	p.Log.Info().Str("page", string(t)).Msg("Added new page")
}

// RemoveModel removes a model from the Pages instance and logs the action.
func (p *Pages) RemoveModel(t PageType) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	delete(p.Models, t)
	p.Log.Info().Str("page", string(t)).Msg("Removed page")
}

// GetModel retrieves a model from the Pages instance by type.
func (p *Pages) GetModel(t PageType) (tea.Model, bool) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	m, ok := p.Models[t]
	return m, ok
}

// SwitchModel switches to a new model of the given type, if it exists.
func (p *Pages) SwitchModel(t PageType) (tea.Model, tea.Cmd) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if model, ok := p.Models[t]; ok {
		p.Navigation = append(p.Navigation, model)
		return model, model.Init()
	}

	p.Log.Error().Str("page", string(t)).Msg("Failed to switch to page")
	m, err := p.CurrentModel()
	if err != nil {
		p.Log.Error().Err(*err).Msg("Failed to get current model")
		return nil, tea.Quit
	}

	return m, m.Init()
}

// SwitchToPreviousModel switches back to the previous model in the navigation stack.
// If the navigation stack is empty, it exits the program.
func (p *Pages) SwitchToPreviousModel() (tea.Model, tea.Cmd) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if len(p.Navigation) == 1 {
		return p.Navigation[0], tea.Quit
	}

	p.Navigation = p.Navigation[:len(p.Navigation)-1]
	return p.Navigation[len(p.Navigation)-1], nil
}

// CurrentModel returns the current model from the navigation stack.
func (p *Pages) CurrentModel() (tea.Model, *error) {
	if len(p.Navigation) == 0 {
		err := errors.New("no models in navigation stack")
		return nil, &err
	}
	return p.Navigation[len(p.Navigation)-1], nil
}

// Theme and styling variables using lipgloss.
var (
	ThemeColorPrimary   = "#aae639"
	ThemeColorSecondary = "#f2ffd8"
	TextPrimaryStyle    = lipgloss.Color(ThemeColorPrimary)
	TextSecondaryStyle  = lipgloss.Color(ThemeColorSecondary)
)
