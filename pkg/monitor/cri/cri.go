package cri

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spacelavr/monitor/pkg/log"
)

// Cri
type Cri struct {
	cli *client.Client
	ctx context.Context
}

// New returns new runtime
func New() *Cri {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	return &Cri{
		cli,
		ctx,
	}
}

// Ping check docker connection
func (r *Cri) Ping() error {
	_, err := r.cli.Ping(r.ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// ContainerInspect returns container info
func (r *Cri) ContainerInspect(id string) (*types.ContainerJSON, error) {
	container, err := r.cli.ContainerInspect(r.ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &container, nil
}

// ContainerList returns a list of running containers
func (r *Cri) ContainerList() ([]types.Container, error) {
	containers, err := r.cli.ContainerList(r.ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return containers, nil
}

// ContainerStats returns container metrics channel
func (r *Cri) ContainerStats(id string) (io.ReadCloser, error) {
	cStats, err := r.cli.ContainerStats(r.ctx, id, true)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return cStats.Body, err
}
