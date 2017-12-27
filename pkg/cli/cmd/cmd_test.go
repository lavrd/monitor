package cmd_test

import (
	"encoding/json"
	"github.com/docker/cli/cli/command/formatter"
	"github.com/lavrs/dlm/pkg/cli/cmd"
	"github.com/lavrs/dlm/pkg/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const container = "container"

func TestGetContainerLogs(t *testing.T) {
	const testLogs = "test logs"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs, err := json.Marshal(map[string]string{
			"logs": testLogs,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(logs)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	cmd.ContainerLogsCmd(container)

	logs, err := cmd.GetContainerLogs(container)
	assert.NoError(t, err)
	assert.Equal(t, testLogs, logs)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	context.Get().Address = ts.URL

	cmd.ContainerLogsCmd(container)
}

func TestGetStoppedContainers(t *testing.T) {
	var testStopped = []string{"container1", "container2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string][]string{
			"stopped": testStopped,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(stopped)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	cmd.StoppedContainersCmd()

	stopped, err := cmd.GetStoppedContainers()
	assert.NoError(t, err)
	assert.Equal(t, testStopped, stopped)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string]interface{}{
			"stopped": nil,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(stopped)
		assert.NoError(t, err)
	}))

	context.Get().Address = ts.URL

	cmd.StoppedContainersCmd()

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	context.Get().Address = ts.URL

	cmd.StoppedContainersCmd()
}

func TestGetLaunchedContainers(t *testing.T) {
	var testLaunched = []string{"container1", "container2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		launched, err := json.Marshal(map[string][]string{
			"launched": testLaunched,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(launched)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	cmd.LaunchedContainersCmd()

	launched, err := cmd.GetLaunchedContainers()
	assert.NoError(t, err)
	assert.Equal(t, testLaunched, launched)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		launched, err := json.Marshal(map[string]interface{}{
			"launched": nil,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(launched)
		assert.NoError(t, err)
	}))

	context.Get().Address = ts.URL

	cmd.LaunchedContainersCmd()

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	context.Get().Address = ts.URL

	cmd.LaunchedContainersCmd()
}

func TestGetContainersMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cs    formatter.ContainerStats
			stats []*formatter.ContainerStats
		)

		cs.SetStatistics(formatter.StatsEntry{
			ID: container,
		})
		stats = append(stats, &cs)

		metrics, err := json.Marshal(map[string]*[]*formatter.ContainerStats{
			"metrics": &stats,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(metrics)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	cmd.ContainersMetricsCmd([]string{"all"})

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics, err := json.Marshal(map[string]*[]formatter.ContainerStats{
			"metrics": nil,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(metrics)
		assert.NoError(t, err)
	}))

	context.Get().Address = ts.URL

	cmd.ContainersMetricsCmd([]string{"all"})

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	context.Get().Address = ts.URL
	cmd.ContainersMetricsCmd([]string{"all"})
}

func TestAPIStatusCmd(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	cmd.APIStatusCmd()

	err := cmd.GetAPIStatus()
	assert.NoError(t, err)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	context.Get().Address = ts.URL

	cmd.APIStatusCmd()

	err = cmd.GetAPIStatus()
	assert.Error(t, err)
}
