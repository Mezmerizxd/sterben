package tui

import (
	"sterben/pkg/log"
	"sterben/pkg/pages"
	"sterben/tui/youtube"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	Home pages.PageType = pages.PageType{
		ID:   "home",
		Name: "Home",
	}
	Youtube pages.PageType = pages.PageType{
		ID:   "youtube",
		Name: "Youtube",
	}
	ImageToIcon pages.PageType = pages.PageType{
		ID:   "image_to_icon",
		Name: "Image to Icon",
	}
)

type Tui struct {
	Log   log.Log
	Pages *pages.Pages
}

func Initialize() (*Tui, error) {
	// Initialize Log
	l := log.New(log.Config{
		Feature:       "tui",
		ConsoleOutput: false,
		FileOutput:    true,
	})
	defer l.Close()

	// Initialize Pages
	p := pages.Initialize(pages.Config{
		Log: l,
	})
	// defer p.Close()

	// Home Page
	homePage := HomePage(&pages.ModelConfig{
		Log:   l,
		Pages: p,
	})
	imageToIconPage := ImageToIconPage(&pages.ModelConfig{
		Log:   l,
		Pages: p,
	})

	// Youtube Home Page
	youtubePage := youtube.HomePage(&pages.ModelConfig{
		Log:   l,
		Pages: p,
	})
	// Youtube set url page
	setUrlPage := youtube.SetUrlPage(&pages.ModelConfig{
		Log:   l,
		Pages: p,
	})

	p.AddModel(Home, homePage)
	p.AddModel(ImageToIcon, imageToIconPage)
	p.AddModel(youtube.Home, youtubePage)
	p.AddModel(youtube.SetUrl, setUrlPage)

	return &Tui{
		Log:   l,
		Pages: p,
	}, nil
}

func (y *Tui) Start() error {
	y.Pages.SwitchModel(Home)

	tea := tea.NewProgram(y.Pages, tea.WithAltScreen(), tea.WithMouseAllMotion())

	_, err := tea.Run()
	if err != nil {
		panic(err)
	}

	return nil
}

func (y *Tui) Close() {
	y.Pages.Close()
	y.Log.Close()
}
