package image_convert

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	ico "github.com/Kodeworks/golang-image-ico"
	"golang.org/x/image/draw"
)

// SupportedImageType represents the type of image supported
type SupportedImageType int

const (
	PNG SupportedImageType = iota
	JPEG
	WEBP
)

var Sizes []image.Point = []image.Point{
	{16, 16},
	{24, 24},
	{32, 32},
	{48, 48},
	{64, 64},
	{128, 128},
	{256, 256},
}

func ConvertImageToMultipleIcons(inputPath string, sizes []image.Point) error {
	// Expand the inputPath to handle the tilde (~) if present
	expandedInputPath, err := expandPath(inputPath)
	if err != nil {
		return fmt.Errorf("error expanding input path: %w", err)
	}

	// Open the input image file
	file, err := os.Open(expandedInputPath)
	if err != nil {
		return fmt.Errorf("error opening image: %w", err)
	}
	defer file.Close()

	var fileName string = "./" + filepath.Base(file.Name())
	// remove extension
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Determine the image type
	imgType, err := getImageType(file)
	if err != nil {
		return fmt.Errorf("error determining image type: %w", err)
	}

	// Seek to the beginning of the file after reading the header
	file.Seek(0, 0)

	// Decode the image
	var img image.Image
	switch imgType {
	case PNG:
		img, err = png.Decode(file)
	case JPEG:
		img, err = jpeg.Decode(file)
	// Add WEBP decoding if needed
	default:
		return fmt.Errorf("unsupported image format")
	}
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	// Create the output directory if it doesn't exist
	if err := makeDirectoryIfNotExists(fileName); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	// Create ICO files for each size
	for _, size := range sizes {
		resizedImg := resizeImage(img, size)
		rgbaImg := convertToRGBA(resizedImg)

		// Prepare the output ICO file path
		outputPath := filepath.Join(fileName, fmt.Sprintf("%dx%d.ico", size.X, size.Y))
		icoFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating ICO file for size %dx%d: %w", size.X, size.Y, err)
		}

		// Encode the resized image as an ICO file
		if err := ico.Encode(icoFile, rgbaImg); err != nil {
			icoFile.Close()
			return fmt.Errorf("error encoding ICO for size %dx%d: %w", size.X, size.Y, err)
		}

		icoFile.Close()
	}

	return nil
}

func ConvertImageToIcon(inputPath string, size image.Point) error {
	return ConvertImageToMultipleIcons(inputPath, []image.Point{size})
}

// convertToRGBA converts the image to RGBA format
func convertToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

// resizeImage resizes the image to the specified size
func resizeImage(img image.Image, size image.Point) image.Image {
	resizedImg := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Src, nil)
	return resizedImg
}

// getImageType determines the image type based on the file signature
func getImageType(file *os.File) (SupportedImageType, error) {
	header := make([]byte, 512)
	_, err := file.Read(header)
	if err != nil {
		return -1, err
	}

	if header[1] == 'P' && header[2] == 'N' && header[3] == 'G' {
		return PNG, nil
	} else if header[6] == 'J' && header[7] == 'F' && header[8] == 'I' && header[9] == 'F' {
		return JPEG, nil
	} else if header[8] == 'W' && header[9] == 'E' && header[10] == 'B' && header[11] == 'P' {
		return WEBP, nil
	} else {
		return -1, fmt.Errorf("unsupported image type")
	}
}

// makeDirectoryIfNotExists creates a directory if it doesn't exist
func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Expand the tilde (~) in the input path to the user's home directory
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error retrieving home directory: %w", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}
