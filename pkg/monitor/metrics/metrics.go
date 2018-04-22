package metrics

import (
	"strings"
	"sync"
	"time"

	"github.com/spacelavr/monitor/pkg/cri"
	"github.com/spacelavr/monitor/pkg/log"
	"github.com/spf13/viper"
)

type MetricsMap struct {
	sync.Mutex
	Metrics map[string]*cri.ContainerStats
}

type ChangesMap struct {
	sync.Mutex
	Changes map[string]*Change
}

type Change struct {
	// stopped or launched
	Status bool
	// change lifetime
	Lifetime time.Time
}

// type for collect metrics
type Metrics struct {
	cri *cri.Cri

	*MetricsMap
	*ChangesMap
	// check for already collect metrics
	Started bool
	// update containers interval
	CInterval time.Duration
	// update container metrics interval
	CMInterval time.Duration
	// changes map flush interval
	FInterval time.Duration
}

// todo rename
type GeneralInfo struct {
	Metrics  []*cri.ContainerStats `json:"metrics,omitempty"`
	Launched []string              `json:"launched,omitempty"`
	Stopped  []string              `json:"stopped,omitempty"`
	Message  string                `json:"message,omitempty"`
}

func New(r *cri.Cri) *Metrics {
	return &Metrics{
		cri: r,
		MetricsMap: &MetricsMap{
			Metrics: make(map[string]*cri.ContainerStats),
		},
		ChangesMap: &ChangesMap{
			Changes: make(map[string]*Change),
		},
		CInterval:  time.Second * time.Duration(viper.GetInt("CInterval")),
		CMInterval: time.Second * time.Duration(viper.GetInt("CMInterval")),
		FInterval:  time.Second * time.Duration(viper.GetInt("FInterval")),
	}
}

// Collect collect metrics (check for new containers)
func (m *Metrics) Collect() error {
	// if metrics already collects, returns
	if m.Started {
		return nil
	}
	m.Started = true

	// check for new containers
	for range time.Tick(m.CInterval) {
		containers, err := m.cri.ContainerList()
		if err != nil {
			return err
		}

		for _, container := range containers {
			if _, ok := m.Metrics[container.Names[0][1:]]; !ok {
				// start collect new container metrics
				go func() {
					if err := m.collect(container.Names[0][1:]); err != nil {
						log.Fatal(err)
					}
				}()
				log.Debug("new container `", container.Names[0][1:], "`")
			}
		}
	}

	return nil
}

// Get returns general info about containers metrics
func (m *Metrics) Get(id string) *GeneralInfo {
	var (
		metrics                []*cri.ContainerStats
		ids, launched, stopped []string
		isNotExist             = 0
	)

	// parse id (all / one / ... containers)
	if id == "all" {
		for _, d := range m.Metrics {
			ids = append(ids, d.Name)
		}
	} else if strings.Contains(id, " ") {
		ids = strings.Split(strings.Replace(id, " ", "", -1), " ")
	} else {
		ids = append(ids, id)
	}
	log.Debug("get", ids, "containers metrics")

	// parse changes
	launched, stopped = m.parseChanges()

	// return if no running containers
	if len(m.Metrics) == 0 {
		log.Debug("no running containers")
		return &GeneralInfo{
			Launched: launched,
			Stopped:  stopped,
			Message:  "no running containers",
		}
	}

	// get containers metrics from data map
	for _, id := range ids {
		if data, ok := m.Metrics[id]; ok {
			metrics = append(metrics, data)
			continue
		}
		// if container are not running
		log.Debug("container `", id, "` are not running")
		isNotExist++
	}
	// returns if all specified containers are not running
	if isNotExist == len(ids) {
		log.Debug("these containers", ids, "are not running")
		return &GeneralInfo{
			Launched: launched,
			Stopped:  stopped,
			Message:  "these containers are not running",
		}
	}

	// returns metrics
	log.Debug(ids, "metrics", metrics)
	return &GeneralInfo{
		Metrics:  metrics,
		Launched: launched,
		Stopped:  stopped,
	}
}
