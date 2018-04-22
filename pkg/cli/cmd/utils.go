package cmd

//
// import (
// 	"encoding/json"
//
// 	"github.com/docker/cli/cli/command/formatter"
// 	h "github.com/spacelavr/monitor/pkg/cli/http"
// 	"github.com/spacelavr/monitor/pkg/context"
// 	"github.com/spacelavr/monitor/pkg/monitor/metrics"
// 	"github.com/spacelavr/monitor/pkg/logger"
// )
//
// const api = "/api/"
//
// // GetContainersMetrics returns containers metrics
// func GetContainersMetrics(id string) ([]formatter.ContainerStats, error) {
// 	body, err := h.GET(context.Get().Address + api + "metrics/" + id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var api metrics.API
// 	if err = json.Unmarshal(body, &api); err != nil {
// 		return nil, err
// 	}
//
// 	logger.Info(id, "containers metrics", api.Metrics)
// 	return api.Metrics, nil
// }
//
// // GetStoppedContainers returns stopped containers
// func GetStoppedContainers() ([]string, error) {
// 	body, err := h.GET(context.Get().Address + api + "stopped")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var api metrics.API
// 	if err = json.Unmarshal(body, &api); err != nil {
// 		return nil, err
// 	}
//
// 	logger.Info("stopped containers", api.Stopped)
// 	return api.Stopped, nil
// }
//
// // GetLaunchedContainers returns launched containers
// func GetLaunchedContainers() ([]string, error) {
// 	body, err := h.GET(context.Get().Address + api + "launched")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var api metrics.API
// 	if err = json.Unmarshal(body, &api); err != nil {
// 		return nil, err
// 	}
//
// 	logger.Info("launched containers", api.Launched)
// 	return api.Launched, nil
// }
//
// // GetAPIStatus returns API status
// func GetAPIStatus() error {
// 	_, err := h.GET(context.Get().Address + api + "status")
// 	return err
// }
