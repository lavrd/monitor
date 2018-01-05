package metrics

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/spacelavr/dlm/pkg/kit/docker"
	"github.com/spacelavr/dlm/pkg/logger"
)

// collect metrics (container stats)
func (m *metrics) collect(id string) {
	fmt.Println("MAGIC")
	// set that container launched
	m.changesMap.Lock()
	m.changes[id] = struct {
		status   bool
		lifetime time.Time
	}{status: true, lifetime: time.Now()}
	m.changesMap.Unlock()
	defer func() {
		logger.Info("container `", id, "` stopped")
		m.setContainerStoppedState(id)
	}()

	// check for container stopped
	var stopped = make(chan bool)
	defer func() {
		fmt.Println("WHIT THIS?")
		close(stopped)
	}()
	go func() {
		for range time.Tick(time.Second) {
			if info, err := docker.ContainerInspect(id); err != nil || !info.State.Running {
				stopped <- true
				return
			}
		}
	}()

	// get container stats channel
	reader, err := docker.ContainerStats(id)
	if err != nil {
		logger.Panic(err)
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var statsJSON *types.StatsJSON

	for range time.Tick(m.updCMetricsInterval) {
		select {
		// container stopped
		case <-stopped:
			fmt.Println("ALARM")
			return
		default:
			// parse metrics
			if err = dec.Decode(&statsJSON); err != nil {
				fmt.Println("KEK")
				return
			}
			// formatting metrics
			metrics := docker.Formatting(id, statsJSON)

			// update metrics
			m.metricsMap.Lock()
			m.metrics[id] = metrics
			m.metricsMap.Unlock()
		}
	}
}

// remove container from metrics map
// and set that container stopped
func (m *metrics) setContainerStoppedState(id string) {
	// set that container stopped
	m.changesMap.Lock()
	m.changes[id] = struct {
		status   bool
		lifetime time.Time
	}{status: false, lifetime: time.Now()}
	m.changesMap.Unlock()

	// remove from metrics map
	m.metricsMap.Lock()
	delete(m.metrics, id)
	m.metricsMap.Unlock()
}

func (m *metrics) parseChanges() (launched, stopped []string) {
	// parse changes
	if len(m.changes) != 0 {
		for id, status := range m.changes {
			if status.status {
				launched = append(launched, id)
			} else {
				stopped = append(stopped, id)
			}
		}

		// flush changes
		for k, v := range m.changes {
			// if the time has passed more than m.flushInterval
			if time.Since(v.lifetime) > m.chFlushInterval {
				delete(m.changes, k)
			}
		}
	}

	logger.Info("stopped", stopped, "launched", launched)
	return
}
