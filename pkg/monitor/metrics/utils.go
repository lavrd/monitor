package metrics

import (
	"encoding/json"
	"strings"
	"time"

	dt "github.com/docker/docker/api/types"
	"github.com/spacelavr/monitor/pkg/cri"
	"github.com/spacelavr/monitor/pkg/types"
	"github.com/spacelavr/monitor/pkg/utils/log"
	"github.com/spf13/viper"
)

// collect collect metrics
func (m *Metrics) collect(id string) error {
	reader, err := m.Cri.ContainerStats(id)
	if err != nil {
		return err
	}
	defer reader.Close()

	var (
		stopped   = make(chan bool)
		statsJSON *dt.StatsJSON
		dec       = json.NewDecoder(reader)
	)

	defer func() {
		close(stopped)
		m.delete(id)
		log.Debug("container `", id, "` stopped")
	}()

	go func() {
		for range time.Tick(time.Second * time.Duration(viper.GetInt(types.FCInterval))) {
			if info, err := m.Cri.ContainerInspect(id); err != nil || !info.State.Running {
				stopped <- true
				return
			}
		}
	}()

	for range time.Tick(m.cmInterval) {
		select {
		case <-stopped:
			return nil
		default:
			if err = dec.Decode(&statsJSON); err != nil {
				return err
			}
			m.save(id, m.Cri.Formatting(id, statsJSON))
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
		if data, ok := m.metrics[id]; ok {
			metrics = append(metrics, data)
			continue
		}
	}

	return metrics
}

// load load metrics from map by id
func (m *metricsMap) load(id string) *cri.ContainerStats {
	m.RLock()
	defer m.RUnlock()

	return m.metrics[id]
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
