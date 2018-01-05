package metrics

import (
	"sync"
	"time"

	"github.com/docker/cli/cli/command/formatter"
)

type metricsMap struct {
	sync.Mutex
	metrics map[string]*formatter.ContainerStats
}

type changesMap struct {
	sync.Mutex
	changes map[string]change
}

type change struct {
	// stopped or launched
	status bool
	// change lifetime
	lifetime time.Time
}

// type for collect metrics
type metrics struct {
	metricsMap
	changesMap
	// check for already collect metrics
	started bool
	// update containers interval
	updContsInterval time.Duration
	// update container metrics interval
	updCMetricsInterval time.Duration
	// changes map flush interval
	chFlushInterval time.Duration
}

// API type for
type API struct {
	Metrics  []formatter.ContainerStats `json:"metrics,omitempty"`
	Launched []string                   `json:"launched,omitempty"`
	Stopped  []string                   `json:"stopped,omitempty"`
	Logs     string                     `json:"logs,omitempty"`
	Message  string                     `json:"message,omitempty"`
}
