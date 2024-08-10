package system_info

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemHostInfo struct {
	Hostname        string `json:"hostname"`
	Uptime          uint64 `json:"uptime"`
	BootTime        uint64 `json:"bootTime"`
	Procs           uint64 `json:"procs"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platformFamily"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
	KernelArch      string `json:"kernelArch"`
	// VirtualizationSystem string `json:"virtualizationSystem"`
	// VirtualizationRole   string `json:"virtualizationRole"`
	// HostID               string `json:"hostid"`
}

type SystemCPUInfo struct {
	ModelName    string `json:"modelName"`
	CoresCount   int32  `json:"cores_count"`
	ThreadsCount int32  `json:"threads_count"`
	Cores        []struct {
		CPU        int32   `json:"cpu"`
		VendorID   string  `json:"vendorId"`
		Family     string  `json:"family"`
		Model      string  `json:"model"`
		Stepping   int32   `json:"stepping"`
		PhysicalID string  `json:"physicalId"`
		CoreID     string  `json:"coreId"`
		Cores      int32   `json:"cores"`
		ModelName  string  `json:"modelName"`
		Mhz        float64 `json:"mhz"`
		CacheSize  int32   `json:"cacheSize"`
		// Flags      []string `json:"flags"`
		// Microcode  string   `json:"microcode"`
	} `json:"cores"`
}

type SystemMemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

type SystemInfo struct {
	Host   SystemHostInfo   `json:"host"`
	CPU    SystemCPUInfo    `json:"cpu"`
	Memory SystemMemoryInfo `json:"memory"`
}

func GetSystemInfo() (SystemInfo, error) {
	hostInfo, err := GetSystemHostInfo()
	if err != nil {
		return SystemInfo{}, err
	}

	cpuInfo, err := GetSystemCPUInfo()
	if err != nil {
		return SystemInfo{}, err
	}

	memoryInfo, err := GetSystemMemoryInfo()
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		Host:   hostInfo,
		CPU:    cpuInfo,
		Memory: memoryInfo,
	}, nil
}

func GetSystemHostInfo() (SystemHostInfo, error) {
	osInfo, err := host.Info()
	if err != nil {
		return SystemHostInfo{}, err
	}

	return SystemHostInfo{
		Hostname:        osInfo.Hostname,
		Uptime:          osInfo.Uptime,
		BootTime:        osInfo.BootTime,
		Procs:           osInfo.Procs,
		OS:              osInfo.OS,
		Platform:        osInfo.Platform,
		PlatformFamily:  osInfo.PlatformFamily,
		PlatformVersion: osInfo.PlatformVersion,
		KernelVersion:   osInfo.KernelVersion,
		KernelArch:      osInfo.KernelArch,
	}, nil
}

func GetSystemCPUInfo() (SystemCPUInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return SystemCPUInfo{}, err
	}

	cpu := SystemCPUInfo{
		ModelName:    cpuInfo[0].ModelName,
		CoresCount:   int32(len(cpuInfo)),
		ThreadsCount: int32(cpuInfo[0].Cores),
	}

	for _, core := range cpuInfo {
		cpu.Cores = append(cpu.Cores, struct {
			CPU        int32   `json:"cpu"`
			VendorID   string  `json:"vendorId"`
			Family     string  `json:"family"`
			Model      string  `json:"model"`
			Stepping   int32   `json:"stepping"`
			PhysicalID string  `json:"physicalId"`
			CoreID     string  `json:"coreId"`
			Cores      int32   `json:"cores"`
			ModelName  string  `json:"modelName"`
			Mhz        float64 `json:"mhz"`
			CacheSize  int32   `json:"cacheSize"`
		}{
			CPU:        core.CPU,
			VendorID:   core.VendorID,
			Family:     core.Family,
			Model:      core.Model,
			Stepping:   core.Stepping,
			PhysicalID: core.PhysicalID,
			CoreID:     core.CoreID,
			Cores:      core.Cores,
			ModelName:  core.ModelName,
			Mhz:        core.Mhz,
			CacheSize:  core.CacheSize,
		})
	}

	return cpu, nil
}

func GetSystemMemoryInfo() (SystemMemoryInfo, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return SystemMemoryInfo{}, err
	}

	return SystemMemoryInfo{
		Total:       memInfo.Total,
		Available:   memInfo.Available,
		Used:        memInfo.Used,
		UsedPercent: memInfo.UsedPercent,
		Free:        memInfo.Free,
	}, nil
}
