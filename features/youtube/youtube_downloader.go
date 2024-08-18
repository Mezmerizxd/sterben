package youtube

import (
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
)

// VideoMetaData holds metadata information for a YouTube video.
type VideoMetaData struct {
	Title       string `json:"title"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	ViewCount   int    `json:"view_count"`
}

var (
	command            = "yt-dlp" // Command to execute yt-dlp.
	downloadWindowsExe = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	downloadUnixBinary = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
)

// DownloadYoutubeVideo downloads a YouTube video to the specified output directory.
func DownloadYoutubeVideo(url, outputDir string) error {
	if !CheckIfYtdlpInstalled() {
		return errors.New("yt-dlp is not installed")
	}

	// Build the command with output directory and URL.
	cmd := exec.Command(command, "-o", filepath.Join(outputDir, "%(title)s.%(ext)s"), url)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// GetVideoMetaData retrieves metadata for the specified YouTube video URL.
func GetVideoMetaData(url string) (*VideoMetaData, error) {
	if !CheckIfYtdlpInstalled() {
		return nil, errors.New("yt-dlp is not installed")
	}

	// Get video metadata in JSON format.
	cmd := exec.Command(command, "-j", url)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON output into the VideoMetaData struct.
	var metadata VideoMetaData
	if err = json.Unmarshal(out, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// CheckIfYtdlpInstalled checks if yt-dlp is installed and available.
func CheckIfYtdlpInstalled() bool {
	// Check if yt-dlp is in the system's PATH.
	if err := exec.Command(command, "--version").Run(); err == nil {
		return true
	}

	// Attempt to find yt-dlp executable in the current directory (Windows).
	if runtime.GOOS == "windows" {
		exePath, err := exec.LookPath("./yt-dlp.exe")
		if err == nil {
			command = exePath
			return true
		}
	}

	// Attempt to find yt-dlp binary in the current directory (Unix).
	if exePath, err := exec.LookPath("./yt-dlp"); err == nil {
		command = exePath
		return true
	}

	return false
}

// DownloadYtdlp downloads the appropriate yt-dlp executable for the current OS.
func DownloadYtdlp() error {
	switch runtime.GOOS {
	case "windows":
		if err := downloadFile(downloadWindowsExe, "yt-dlp.exe"); err != nil {
			return err
		}
		command = "yt-dlp.exe"
	case "linux", "darwin":
		if err := downloadFile(downloadUnixBinary, "yt-dlp"); err != nil {
			return err
		}
		if err := exec.Command("chmod", "+x", "yt-dlp").Run(); err != nil {
			return err
		}
		command = "./yt-dlp"
	default:
		return errors.New("unsupported OS")
	}

	return nil
}

// downloadFile downloads a file from the specified URL to the given filename.
func downloadFile(url, filename string) error {
	cmd := exec.Command("curl", "-L", "-o", filename, url)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
