package metrics

import (
	"encoding/json"
	"strings"
	"time"

	"monitor/pkg/cri"
	"monitor/pkg/utils/log"

	"github.com/docker/docker/api/types"
)

// collect collect metrics
func (m *Metrics) collect(id string) error {
	defer func() {
		m.delete(id)
		log.Debug("container `", id, "` stopped")
	}()

	reader, err := m.Cri.ContainerStats(id)
	if err != nil {
		return err
	}
	defer reader.Close()

	var (
		stats *types.StatsJSON
		dec   = json.NewDecoder(reader)
	)

	for range time.Tick(m.cmInterval) {
		if err = dec.Decode(&stats); err != nil {
			return nil
		} else {
			m.save(id, m.Cri.Formatting(id, stats))
		}
	}

	return nil
}

// parse parse ids and returns ids slice
func (m *Metrics) parse(id string) []string {
	var (
		ids []string
	)

	if id == "all" {
		for _, d := range m.metrics {
			ids = append(ids, d.Name)
		}
	} else if strings.Contains(id, " ") {
		ids = strings.Split(strings.Replace(id, " ", "", -1), " ")
	} else {
		ids = append(ids, id)
	}

	return ids
}

// accumulate accumulate metrics by ids
func (m *Metrics) accumulate(ids []string) []*cri.ContainerStats {
	var (
		metrics []*cri.ContainerStats
	)

	for _, id := range ids {
		if data, ok := m.load(id); ok {
			metrics = append(metrics, data)
		}
	}

	return metrics
}

// load load metrics from map by id
func (m *metricsMap) load(id string) (*cri.ContainerStats, bool) {
	m.RLock()
	defer m.RUnlock()

	cs, ok := m.metrics[id]
	return cs, ok
}

// save save metrics to map by id
func (m *metricsMap) save(id string, metrics *cri.ContainerStats) {
	m.Lock()
	defer m.Unlock()

	m.metrics[id] = metrics
}

// delete delete metrics from map by id
func (m *metricsMap) delete(id string) {
	m.Lock()
	defer m.Unlock()

	delete(m.metrics, id)
}
