package docker

import (
	"math"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

// ContainerStats
type ContainerStats struct {
	Name             string    `json:"name"`
	Id               string    `json:"id"`
	CPUPercentage    float64   `json:"cpu_percentage"`
	Memory           float64   `json:"memory"`
	MemoryLimit      float64   `json:"memory_limit"`
	MemoryPercentage float64   `json:"memory_percentage"`
	NetworkRx        float64   `json:"network_rx"`
	NetworkTx        float64   `json:"network_tx"`
	IORead           float64   `json:"io_read"`
	IOWrite          float64   `json:"io_write"`
	Time             time.Time `json:"time"`
}

// Formatting returns formatted container stats
func (d *Docker) Formatting(id string, s *types.StatsJSON) *ContainerStats {
	cs := &ContainerStats{
		Name:        id,
		Id:          s.ID,
		Memory:      float64(s.MemoryStats.Usage),
		MemoryLimit: float64(s.MemoryStats.Limit),
		Time:        time.Now().UTC(),
	}

	cs.memory(s)
	cs.io(s.BlkioStats)
	cs.network(s.Networks)
	cs.cpu(s)

	return cs
}

// memory set memory in percentages
func (cs *ContainerStats) memory(statsJSON *types.StatsJSON) {
	mem := float64(statsJSON.MemoryStats.Usage) / float64(statsJSON.MemoryStats.Limit) * 100.0

	if math.IsNaN(mem) {
		cs.MemoryPercentage = 0
	}

	cs.MemoryPercentage = mem
}

// network set network metrics
func (cs *ContainerStats) network(network map[string]types.NetworkStats) {
	var (
		rx, tx float64
	)

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}

	cs.NetworkTx = tx
	cs.NetworkRx = rx
}

// cpu set cpu in percentages
func (cs *ContainerStats) cpu(s *types.StatsJSON) {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta = float64(s.CPUStats.SystemUsage) - float64(s.PreCPUStats.SystemUsage)
		onlineCPUs  = float64(s.CPUStats.CPUUsage.TotalUsage)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(s.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	cs.CPUPercentage = cpuPercent
}

// io set r/w metrics
func (cs *ContainerStats) io(IOStats types.BlkioStats) {
	var (
		blkRead, blkWrite uint64
	)

	for _, io := range IOStats.IoServiceBytesRecursive {
		switch strings.ToLower(io.Op) {
		case "read":
			blkRead = blkRead + io.Value
		case "write":
			blkWrite = blkWrite + io.Value
		}
	}

	cs.IORead = float64(blkRead)
	cs.IOWrite = float64(blkWrite)
}
