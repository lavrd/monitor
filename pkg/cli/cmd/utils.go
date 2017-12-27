package cmd

import (
	"encoding/json"
	"github.com/docker/cli/cli/command/formatter"
	h "github.com/lavrs/dlm/pkg/cli/http"
	"github.com/lavrs/dlm/pkg/context"
	"github.com/lavrs/dlm/pkg/kit/metrics"
	"github.com/lavrs/dlm/pkg/logger"
)

const api = "/api/"

// GetContainerLogs returns container logs
func GetContainerLogs(id string) (string, error) {
	body, err := h.GET(context.Get().Address + api + "logs/" + id)
	if err != nil {
		return "", err
	}

	var api metrics.API
	if err = json.Unmarshal(body, &api); err != nil {
		return "", err
	}

	logger.Info(id, "container logs", api.Logs)
	return api.Logs, nil
}

// GetContainersMetrics returns containers metrics
func GetContainersMetrics(id string) ([]formatter.ContainerStats, error) {
	body, err := h.GET(context.Get().Address + api + "metrics/" + id)
	if err != nil {
		return nil, err
	}

	var api metrics.API
	if err = json.Unmarshal(body, &api); err != nil {
		return nil, err
	}

	logger.Info(id, "containers metrics", api.Metrics)
	return api.Metrics, nil
}

// GetStoppedContainers returns stopped containers
func GetStoppedContainers() ([]string, error) {
	body, err := h.GET(context.Get().Address + api + "stopped")
	if err != nil {
		return nil, err
	}

	var api metrics.API
	if err = json.Unmarshal(body, &api); err != nil {
		return nil, err
	}

	logger.Info("stopped containers", api.Stopped)
	return api.Stopped, nil
}

// GetLaunchedContainers returns launched containers
func GetLaunchedContainers() ([]string, error) {
	body, err := h.GET(context.Get().Address + api + "launched")
	if err != nil {
		return nil, err
	}

	var api metrics.API
	if err = json.Unmarshal(body, &api); err != nil {
		return nil, err
	}

	logger.Info("launched containers", api.Launched)
	return api.Launched, nil
}

// GetAPIStatus returns API status
func GetAPIStatus() error {
	_, err := h.GET(context.Get().Address + "/status")
	return err
}
