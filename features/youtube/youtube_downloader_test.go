package youtube

import (
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	err := DownloadYoutubeVideo("https://www.youtube.com/watch?v=Tkb2yVr8kfY", "./output")
	if err != nil {
		t.Error("Error:", err)
	} else {
		// Delete the downloaded folder
		err := os.RemoveAll("./output")
		if err != nil {
			t.Error("Error:", err)
		}

		t.Log("Download successful!")
	}

	meta, err := GetVideoMetaData("https://www.youtube.com/watch?v=Tkb2yVr8kfY")
	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Logf("Title: %s\nID: %s\nDescription: %s\nDuration: %d\nView Count: %d\n", meta.Title, meta.ID, meta.Description, meta.Duration, meta.ViewCount)
	}
}
