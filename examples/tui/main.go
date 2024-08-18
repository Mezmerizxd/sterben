package main

import (
	"fmt"
	"sterben/features/youtube"
	"sterben/tui"
)

func main() {
	if !youtube.CheckIfYtdlpInstalled() {
		fmt.Println("yt-dlp is not installed, installing...	")
		// Attempt to install yt-dlp
		err := youtube.DownloadYtdlp()
		if err != nil {
			panic(err)
		}
	}

	y, err := tui.Initialize()
	if err != nil {
		panic(err)
	}

	err = y.Start()
	if err != nil {
		panic(err)
	}
}
