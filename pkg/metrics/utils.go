package metrics

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/rs/zerolog/log"

	"monitor/pkg/docker"
)

// collect collect metrics
func (m *Metrics) collect(id string) error {
	defer func() {
		m.delete(id)
		log.Debug().Msgf("container %s stopped", id)
	}()

	reader, err := m.docker.ContainerStats(id)
	if err != nil {
		return err
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Error().Err(err).Msg("close reader error")
		}
	}()

	var stats *types.StatsJSON
	decoder := json.NewDecoder(reader)

	for range time.Tick(m.mInterval) {
		if err = decoder.Decode(&stats); err != nil {
			return nil
		} else {
			m.save(id, m.docker.Formatting(id, stats))
		}
	}

	return nil
}

// parse parse ids and returns ids slice
func (m *Metrics) parse(id string) []string {
	var ids []string

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
func (m *Metrics) accumulate(ids []string) []*docker.ContainerStats {
	var metrics []*docker.ContainerStats

	for _, id := range ids {
		if data, ok := m.load(id); ok {
			metrics = append(metrics, data)
		}
	}

	return metrics
}

// load load metrics from map by id
func (m *metricsMap) load(id string) (*docker.ContainerStats, bool) {
	m.RLock()
	defer m.RUnlock()

	cs, ok := m.metrics[id]
	return cs, ok
}

// save save metrics to map by id
func (m *metricsMap) save(id string, metrics *docker.ContainerStats) {
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
