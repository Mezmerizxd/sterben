package youtube

import (
	"fmt"
	"sterben/features/youtube"
	"sterben/pkg/log"
	"sterben/pkg/pages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	Home pages.PageType = pages.PageType{
		ID:   "youtube_home",
		Name: "Youtube Home",
	}
	SetUrl pages.PageType = pages.PageType{
		ID:   "youtube_set_url",
		Name: "Set Url",
	}
	Download pages.PageType = pages.PageType{
		ID:   "youtube_download",
		Name: "Download",
	}
)

type YoutubeTui struct {
	Log   log.Log
	Pages *pages.Pages
}

func Initialize() (*YoutubeTui, error) {
	// Check if yt-dlp is installed
	if !youtube.CheckIfYtdlpInstalled() {
		fmt.Println("yt-dlp is not installed, installing...	")
		// Attempt to install yt-dlp
		err := youtube.DownloadYtdlp()
		if err != nil {
			return nil, err
		}
	}

	// Initialize Log
	l := log.New(log.Config{
		Feature:       "tui_youtube",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer l.Close()
	l2 := log.New(log.Config{
		Feature:       "tui_youtube_2",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer l2.Close()
	l3 := log.New(log.Config{
		Feature:       "tui_youtube_3",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer l3.Close()

	// Initialize Pages
	p := pages.Initialize(pages.Config{
		Log: l,
	})
	// defer p.Close()

	// Home Page
	homePage := HomePage(&pages.ModelConfig{
		Log:   l2,
		Pages: p,
	})

	// Set Url Page
	setUrlPage := SetUrlPage(&pages.ModelConfig{
		Log:   l3,
		Pages: p,
	})

	p.AddModel(Home, homePage)
	p.AddModel(SetUrl, setUrlPage)

	return &YoutubeTui{
		Log:   l,
		Pages: p,
	}, nil
}

func (y *YoutubeTui) Start() error {
	y.Pages.SwitchModel(Home)

	tea := tea.NewProgram(y.Pages, tea.WithAltScreen(), tea.WithMouseAllMotion())

	_, err := tea.Run()
	if err != nil {
		panic(err)
	}

	return nil
}

func GetDebugData(p *pages.Pages) string {
	var s string

	// Home Data
	homePageModel := p.Models[Home].(*HomePageModel)
	s += "Home Data\n"
	s += "---------\n"
	s += "Current Cursor: " + homePageModel.Options.Cursor.Name + "\n"

	// Set Url Data
	setUrlPageModel := p.Models[SetUrl].(*SetUrlPageModel)
	s += "\nSet Url Data\n"
	s += "---------\n"
	s += "Input: " + setUrlPageModel.Input.Value() + "\n"
	if setUrlPageModel.MetaDataError != "" {
		s += "MetaDataError: " + setUrlPageModel.MetaDataError + "\n"
	}
	if setUrlPageModel.MetaDataLoading {
		s += "MetaDataLoading: true\n"
	}
	if setUrlPageModel.InputError != "" {
		s += "InputError: " + setUrlPageModel.InputError + "\n"
	}
	if setUrlPageModel.MetaData != nil {
		s += "MetaData\n"
		s += "---------\n"
		s += "ID: " + setUrlPageModel.MetaData.ID + "\n"
		s += "Title: " + setUrlPageModel.MetaData.Title + "\n"
		s += fmt.Sprintf("Views: %d\n", setUrlPageModel.MetaData.ViewCount)
		s += fmt.Sprintf("Duration: %d\n", setUrlPageModel.MetaData.Duration)
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("#0cd3ff")).Render(s)
}
