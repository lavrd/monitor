package metrics_test

import (
	//"fmt"
	"testing"
	"time"

	"github.com/spacelavr/dlm/pkg/kit/docker"
	m "github.com/spacelavr/dlm/pkg/kit/metrics"
	"github.com/spacelavr/dlm/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	container   = "splines"
	image       = "bfirsh/reticulate-splines"
	updInterval = time.Millisecond * 500
)

func init() {
	err := docker.ImagePull(image)
	if err != nil {
		logger.Panic(err)
	}
	err = docker.ContainerCreate(image, container)
	if err != nil {
		logger.Panic(err)
	}
	err = docker.ContainerStart(container)
	if err != nil {
		logger.Panic(err)
	}

	m.Get().SetUContsInterval(updInterval)
	m.Get().SetUCMetricsInterval(updInterval)
	m.Get().SetChangesFlushInterval(updInterval)
	go m.Get().Collect()
}

func TestGetMetricsObj(t *testing.T) {
	assert.NotNil(t, m.Get())
}

func TestMetricsAlreadyCollecting(t *testing.T) {
	m.Get().Collect()
}

func TestGetNoRunningContainers(t *testing.T) {
	cMetrics := m.Get().Get(container)
	assert.Equal(t, "no running containers", cMetrics.Message)
}

func TestGetContainerLogs(t *testing.T) {
	pending(updInterval)

	logs := m.GetContainerLogs(container)
	assert.NotEmpty(t, logs)
}

func TestGetLaunchedContainers(t *testing.T) {
	assert.NotEmpty(t, m.Get().GetLaunchedContainers())
}

func TestGet(t *testing.T) {
	pending(updInterval)
	assert.NotNil(t, m.Get().Get("all").Metrics)
}

func TestGetSpecifiedContainers(t *testing.T) {
	assert.Equal(t, "these containers are not running",
		m.Get().Get(container + " container").Message)
}

func TestGetStoppedContainers(t *testing.T) {
	err := docker.ContainerStop(container)
	assert.NoError(t, err)
	pending(updInterval)

	assert.NotEmpty(t, m.Get().GetStoppedContainers())
}

func TestContainerRemoveHandle(t *testing.T) {
	err := docker.ContainerStart(container)
	assert.NoError(t, err)
	pending(updInterval)
	err = docker.ContainerRemove(container)
	assert.NoError(t, err)
	pending(updInterval)
	err = docker.ImageRemove(image)
	assert.NoError(t, err)
}

func pending(t time.Duration) {
	time.Sleep(t * 2)
}
