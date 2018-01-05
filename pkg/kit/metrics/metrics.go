package metrics

import (
	"strings"
	"time"

	"github.com/docker/cli/cli/command/formatter"
	"github.com/spacelavr/dlm/pkg/kit/docker"
	"github.com/spacelavr/dlm/pkg/logger"
)

// init metrics
var m = &metrics{
	metricsMap: metricsMap{
		metrics: make(map[string]*formatter.ContainerStats),
	},
	changesMap: changesMap{
		changes: make(map[string]change),
	},
	updContsInterval:    time.Second * 3,
	updCMetricsInterval: time.Second * 1,
	chFlushInterval:     time.Second * 10,
}

// Get returns metrics obj
func Get() *metrics {
	return m
}

// SetUCMetricsInterval set update container metrics interval
func (m *metrics) SetUCMetricsInterval(t time.Duration) {
	logger.Info("set update container metrics interval", t)
	m.updCMetricsInterval = t
}

// SetUContsInterval set update containers interval
func (m *metrics) SetUContsInterval(t time.Duration) {
	logger.Info("set update containers interval", t)
	m.updContsInterval = t
}

// SetChangesFlushInterval set changes map flush interval
func (m *metrics) SetChangesFlushInterval(t time.Duration) {
	logger.Info("set changes map flush interval", t)
	m.chFlushInterval = t
}

// Collect collect metrics (check for new containers)
func (m *metrics) Collect() {
	// if metrics already collects, returns
	if m.started {
		return
	}
	m.started = true

	// wait for docker connection
	for {
		if err := docker.Ping(); err == nil {
			break
		}
		logger.Info("trying to connect to docker")
		time.Sleep(time.Second)
	}

	// check for new containers
	for range time.Tick(m.updContsInterval) {
		containers, err := docker.ContainerList()
		if err != nil {
			logger.Panic(err)
		}

		for _, container := range containers {
			if _, ok := m.metrics[container.Names[0][1:]]; !ok {
				// start collect new container metrics
				go m.collect(container.Names[0][1:])
				logger.Info("new container `", container.Names[0][1:], "`")
			}
		}
	}
}

// GetContainerLogs returns container logs
func GetContainerLogs(id string) *API {
	return &API{
		Logs: docker.ContainerLogs(id),
	}
}

// GetStoppedContainers returns stopped containers
func (m *metrics) GetStoppedContainers() *API {
	_, stopped := m.parseChanges()
	return &API{
		Stopped: stopped,
	}
}

// GetLaunchedContainers returns Launched containers
func (m *metrics) GetLaunchedContainers() *API {
	launched, _ := m.parseChanges()
	return &API{
		Launched: launched,
	}
}

// Get returns containers metrics
func (m *metrics) Get(id string) *API {
	var (
		metrics                []formatter.ContainerStats
		ids, launched, stopped []string
		isNotExist             = 0
	)

	// parse id (all / one / ... containers)
	if id == "all" {
		for _, d := range m.metrics {
			ids = append(ids, d.Name)
		}
	} else if strings.Contains(id, " ") {
		ids = strings.Split(strings.Replace(id, " ", "", -1), " ")
	} else {
		ids = append(ids, id)
	}
	logger.Info("get", ids, "containers metrics")

	// parse changes
	launched, stopped = m.parseChanges()

	// return if no running containers
	if len(m.metrics) == 0 {
		logger.Info("no running containers")
		return &API{
			Launched: launched,
			Stopped:  stopped,
			Message:  "no running containers",
		}
	}

	// get containers metrics from data map
	for _, id := range ids {
		if data, ok := m.metrics[id]; ok {
			metrics = append(metrics, *data)
			continue
		}
		// if container are not running
		logger.Info("container `", id, "` are not running")
		isNotExist++
	}
	// returns if all specified containers are not running
	if isNotExist == len(ids) {
		logger.Info("these containers", ids, "are not running")
		return &API{
			Launched: launched,
			Stopped:  stopped,
			Message:  "these containers are not running",
		}
	}

	// returns metrics
	logger.Info(ids, "metrics", metrics)
	return &API{
		Metrics:  metrics,
		Launched: launched,
		Stopped:  stopped,
	}
}
