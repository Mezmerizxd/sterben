package pages

import (
	"sterben/pkg/log"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

const testPage PageType = "testPage"

type mockModel struct{}

func (m mockModel) Init() tea.Cmd { return nil }
func (m mockModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m mockModel) View() string { return "mock" }

func TestAddModel(t *testing.T) {
    // Initialize the Pages instance
    p := Initialize(Config{
			Log: log.New(log.Config{
				Feature: "pages_test",
				ConsoleOutput: true,
				FileOutput: false,
			}),
		})

    // Add a mock model
    mock := mockModel{}
    p.AddModel(testPage, mock)

		// Get the model
		model, ok := p.GetModel(testPage)

		if !ok {
			t.Errorf("Failed to get model")
		}

    // Assert that the model was added
    assert.Equal(t, mock, model)
}
