package image_convert

import (
	"testing"
)

func TestImageToIcon(t *testing.T) {
	err := ConvertImageToMultipleIcons("input.png", Sizes)
	if err != nil {
		t.Log("Error:", err)
	} else {
		t.Log("Conversion successful!")
	}
}
