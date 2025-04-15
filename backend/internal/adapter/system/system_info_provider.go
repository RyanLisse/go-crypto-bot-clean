package system

import (
	"context"
	"runtime"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/status"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfoProvider provides system resource information
type SystemInfoProvider struct {
	logger       *zerolog.Logger
	diskPath     string
	lastCPUUsage float64
	lastMemUsage float64
	lastDiskUsage float64
}

// NewSystemInfoProvider creates a new system info provider
func NewSystemInfoProvider(logger *zerolog.Logger, diskPath string) *SystemInfoProvider {
	if diskPath == "" {
		diskPath = "/"
	}
	
	return &SystemInfoProvider{
		logger:   logger,
		diskPath: diskPath,
	}
}

// GetSystemInfo returns the current system resource information
func (p *SystemInfoProvider) GetSystemInfo(ctx context.Context) (*status.SystemInfo, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	// Get CPU usage
	cpuUsage, err := p.getCPUUsage()
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to get CPU usage")
		cpuUsage = p.lastCPUUsage
	} else {
		p.lastCPUUsage = cpuUsage
	}
	
	// Get memory usage
	memUsage, err := p.getMemoryUsage()
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to get memory usage")
		memUsage = p.lastMemUsage
	} else {
		p.lastMemUsage = memUsage
	}
	
	// Get disk usage
	diskUsage, err := p.getDiskUsage()
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to get disk usage")
		diskUsage = p.lastDiskUsage
	} else {
		p.lastDiskUsage = diskUsage
	}
	
	return &status.SystemInfo{
		CPUUsage:            cpuUsage,
		MemoryUsage:         memUsage,
		DiskUsage:           diskUsage,
		NumGoroutines:       runtime.NumGoroutine(),
		AllocatedMemory:     memStats.Alloc,
		TotalAllocatedMemory: memStats.TotalAlloc,
		GCPauseTotal:        memStats.PauseTotalNs,
		LastGCPause:         memStats.PauseNs[(memStats.NumGC+255)%256],
	}, nil
}

// getCPUUsage returns the current CPU usage percentage
func (p *SystemInfoProvider) getCPUUsage() (float64, error) {
	percentage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	
	if len(percentage) == 0 {
		return 0, nil
	}
	
	return percentage[0], nil
}

// getMemoryUsage returns the current memory usage percentage
func (p *SystemInfoProvider) getMemoryUsage() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	
	return memInfo.UsedPercent, nil
}

// getDiskUsage returns the current disk usage percentage
func (p *SystemInfoProvider) getDiskUsage() (float64, error) {
	diskInfo, err := disk.Usage(p.diskPath)
	if err != nil {
		return 0, err
	}
	
	return diskInfo.UsedPercent, nil
}
