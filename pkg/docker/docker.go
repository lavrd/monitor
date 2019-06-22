package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

const CtxTimeout = time.Second * 10

// Docker
type Docker struct {
	client *client.Client
}

// New returns new docker
func New() (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Docker{cli}, nil
}

// Close close docker connection
func (d *Docker) Close() {
	if err := d.client.Close(); err != nil {
		log.Error().Err(err).Msg("close connection with docker error")
	}
}

// ContainerInspect returns container info
func (d *Docker) ContainerInspect(id string) (*types.ContainerJSON, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	container, err := d.client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// ContainerList returns a list of running containers
func (d *Docker) ContainerList() ([]types.Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// ContainerStats returns container metrics channel
func (d *Docker) ContainerStats(id string) (io.ReadCloser, error) {
	cs, err := d.client.ContainerStats(context.Background(), id, true)
	if err != nil {
		return nil, err
	}
	return cs.Body, err
}
