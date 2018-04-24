package metrics

import (
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/spacelavr/monitor/pkg/utils/log"
)

// collect metrics (container stats)
func (m *Metrics) collect(id string) error {
	// set that container launched
	defer func() {
		log.Debug("container `", id, "` stopped")
		m.setContainerStoppedState(id)
	}()

	// check for container stopped
	var stopped = make(chan bool)
	defer func() {
		close(stopped)
	}()
	go func() {
		for range time.Tick(time.Second) {
			if info, err := m.cri.ContainerInspect(id); err != nil || !info.State.Running {
				stopped <- true
				return
			}
		}
	}()

	// get container stats channel
	reader, err := m.cri.ContainerStats(id)
	if err != nil {
		log.Error(err)
		return err
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var statsJSON *types.StatsJSON

	for range time.Tick(m.CMInterval) {
		select {
		// container stopped
		case <-stopped:
			return nil
		default:
			// parse metrics
			if err = dec.Decode(&statsJSON); err != nil {
				return nil
			}

			// formatting metrics
			metrics := m.cri.Formatting(id, statsJSON)

			// update metrics
			m.MetricsMap.Lock()
			m.Metrics[id] = metrics
			m.MetricsMap.Unlock()
		}
	}

	return nil
}

// remove container from metrics map
// and set that container stopped
func (m *Metrics) setContainerStoppedState(id string) {
	// remove from metrics map
	m.MetricsMap.Lock()
	delete(m.Metrics, id)
	m.MetricsMap.Unlock()
}
