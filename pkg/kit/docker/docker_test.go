package docker_test

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/lavrs/dlm/pkg/kit/docker"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	container = "splines"
	image     = "bfirsh/reticulate-splines"
)

func TestImagePull(t *testing.T) {
	err := docker.ImagePull(image)
	assert.NoError(t, err)
}

func TestContainerCreate(t *testing.T) {
	err := docker.ContainerCreate(image, container)
	assert.NoError(t, err)
}

func TestContainerStart(t *testing.T) {
	err := docker.ContainerStart(container)
	assert.NoError(t, err)
}

func TestPing(t *testing.T) {
	err := docker.Ping()
	assert.NoError(t, err)
}

func TestContainerList(t *testing.T) {
	_, err := docker.ContainerList()
	assert.NoError(t, err)
}

func TestContainerInspect(t *testing.T) {
	info, err := docker.ContainerInspect(container)
	assert.NoError(t, err)
	assert.NotNil(t, info)
}

func TestContainerStats(t *testing.T) {
	reader, err := docker.ContainerStats(container)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
}

func TestFormatting(t *testing.T) {
	reader, err := docker.ContainerStats(container)
	assert.NoError(t, err)

	dec := json.NewDecoder(reader)
	var statsJSON *types.StatsJSON
	err = dec.Decode(&statsJSON)
	assert.NoError(t, err)

	assert.NotNil(t, docker.Formatting(container, statsJSON))
}

func TestContainersLogs(t *testing.T) {
	logs := docker.ContainerLogs(container)
	assert.NotEmpty(t, logs)
}

func TestContainerStop(t *testing.T) {
	err := docker.ContainerStop(container)
	assert.NoError(t, err)
}

func TestContainerRemove(t *testing.T) {
	err := docker.ContainerRemove(container)
	assert.NoError(t, err)
}

func TestImageRemove(t *testing.T) {
	err := docker.ImageRemove(image)
	assert.NoError(t, err)
}
