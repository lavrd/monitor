package metrics

import (
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/spacelavr/monitor/pkg/log"
	"github.com/spacelavr/monitor/pkg/monitor/env"
)

// collect metrics (container stats)
func (m *Metrics) collect(id string) error {
	var (
		cri = env.GetCri()
	)

	// set that container launched
	m.ChangesMap.Lock()
	m.Changes[id] = &Change{Status: true, Lifetime: time.Now()}
	m.ChangesMap.Unlock()
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
			if info, err := cri.ContainerInspect(id); err != nil || !info.State.Running {
				stopped <- true
				return
			}
		}
	}()

	// get container stats channel
	reader, err := cri.ContainerStats(id)
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
			metrics := cri.Formatting(id, statsJSON)

			// update metrics
			m.MetricsMap.Lock()
			m.Metrics[id] = metrics
			m.MetricsMap.Unlock()
		}
	}
}

// remove container from metrics map
// and set that container stopped
func (m *Metrics) setContainerStoppedState(id string) {
	// set that container stopped
	m.ChangesMap.Lock()
	m.Changes[id] = &Change{Status: false, Lifetime: time.Now()}
	m.ChangesMap.Unlock()

	// remove from metrics map
	m.MetricsMap.Lock()
	delete(m.Metrics, id)
	m.MetricsMap.Unlock()
}

func (m *Metrics) parseChanges() (launched, stopped []string) {
	// parse changes
	if len(m.Changes) != 0 {
		for id, status := range m.Changes {
			if status.Status {
				launched = append(launched, id)
			} else {
				stopped = append(stopped, id)
			}
		}

		// flush changes
		for k, v := range m.Changes {
			// if the time has passed more than m.flushInterval
			if time.Since(v.Lifetime) > m.FInterval {
				delete(m.Changes, k)
			}
		}
	}

	log.Debug("stopped", stopped, "launched", launched)
	return
}
