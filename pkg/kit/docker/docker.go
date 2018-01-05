package docker

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spacelavr/dlm/pkg/logger"
)

var (
	cli *client.Client
	ctx context.Context
)

func init() {
	var err error

	cli, err = client.NewEnvClient()
	if err != nil {
		logger.Panic(err)
	}

	ctx = context.Background()
}

// Ping check docker connection
func Ping() error {
	_, err := cli.Ping(ctx)
	return err
}

// ContainerInspect returns container info
func ContainerInspect(id string) (types.ContainerJSON, error) {
	return cli.ContainerInspect(ctx, id)
}

// ContainerList returns a list of running containers
func ContainerList() ([]types.Container, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	return containers, err
}

// ContainerStats returns container metrics channel
func ContainerStats(id string) (io.ReadCloser, error) {
	cStats, err := cli.ContainerStats(ctx, id, true)
	return cStats.Body, err
}

// ContainerCreate create container
func ContainerCreate(image, container string) error {
	_, err := cli.ContainerCreate(ctx, &c.Config{
		Image: image,
	}, &c.HostConfig{}, &network.NetworkingConfig{}, container)
	return err
}

// ContainerStart launches container
func ContainerStart(container string) error {
	return cli.ContainerStart(ctx, container, types.ContainerStartOptions{})
}

// ImagePull download image
func ImagePull(image string) error {
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}
	return nil
}

// ContainerStop stops the container
func ContainerStop(container string) error {
	var t = time.Duration(0)
	return cli.ContainerStop(ctx, container, &t)
}

// ImageRemove removes image
func ImageRemove(image string) error {
	_, err := cli.ImageRemove(ctx, image, types.ImageRemoveOptions{
		Force: true,
	})
	return err
}

// ContainerRemove removes container
func ContainerRemove(container string) error {
	return cli.ContainerRemove(ctx, container, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}

// ContainerLogs returns containers logs
func ContainerLogs(container string) string {
	reader, err := cli.ContainerLogs(ctx, container, types.ContainerLogsOptions{
		ShowStdout: true,
	})
	if err != nil {
		return err.Error()
	}
	defer reader.Close()

	logs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err.Error()
	}
	logger.Info(container, "container logs:", logs)
	return string(logs)
}
