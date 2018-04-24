package cri

import (
	"math"
	"strings"

	"github.com/docker/docker/api/types"
)

type ContainerStats struct {
	Name             string  `json:"name"`
	Id               string  `json:"id"`
	CPUPercentage    float64 `json:"cpu_percentage"`
	Memory           float64 `json:"memory"`
	MemoryLimit      float64 `json:"memory_limit"`
	MemoryPercentage float64 `json:"memory_percentage"`
	NetworkRx        float64 `json:"network_rx"`
	NetworkTx        float64 `json:"network_tx"`
	IORead           float64 `json:"io_read"`
	IOWrite          float64 `json:"io_write"`
}

// Formatting returns the basic metrics from all
func (r *Cri) Formatting(id string, s *types.StatsJSON) *ContainerStats {
	cs := &ContainerStats{
		Name:        id,
		Id:          s.ID,
		Memory:      float64(s.MemoryStats.Usage),
		MemoryLimit: float64(s.MemoryStats.Limit),
	}

	cs.memory(s)
	cs.io(s.BlkioStats)
	cs.network(s.Networks)
	cs.cpu(s)

	return cs
}

// parse memory (returns memory in percentages)
func (s *ContainerStats) memory(statsJSON *types.StatsJSON) {
	mem := float64(statsJSON.MemoryStats.Usage) / float64(statsJSON.MemoryStats.Limit) * 100.0

	if math.IsNaN(mem) {
		s.MemoryPercentage = 0
	}

	s.MemoryPercentage = mem
}

// parse network metrics
func (s *ContainerStats) network(network map[string]types.NetworkStats) {
	var (
		rx, tx float64
	)

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}

	s.NetworkTx = tx
	s.NetworkRx = rx
}

// parse cpu (returns cpu in percentages)
func (cs *ContainerStats) cpu(s *types.StatsJSON) {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta = float64(s.CPUStats.SystemUsage) - float64(s.PreCPUStats.SystemUsage)
		onlineCPUs  = float64(s.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(s.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cs.CPUPercentage = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	cs.CPUPercentage = cpuPercent
}

// parse r/w metrics
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
