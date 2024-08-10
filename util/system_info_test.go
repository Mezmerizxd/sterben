package system_info

import (
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	systemInfo, err := GetSystemInfo()
	if err != nil {
		t.Fatalf("Failed to get CPU info: %v", err)
	}

	t.Logf("%+v", systemInfo)
}
