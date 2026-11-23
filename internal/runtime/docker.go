// Package runtime handles a runtime interface and its implementation(s) for containerized applications.
package runtime

import (
	"context"
	"log"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/will-petersmeyer/workload-orchestrator/internal/workload"
)

type DockerRuntime struct {
	apiClient *client.Client
}

func NewDockerRuntime() (*DockerRuntime, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerRuntime{apiClient: apiClient}, nil
}

// Close frees the Docker runtime without disturbing existing containers.
func (d *DockerRuntime) Close() error {
	err := d.apiClient.Close()
	if err != nil {
		return err
	}

	return nil
}

// Create creates a Docker container using the information provided by the given workload Struct.
func (d *DockerRuntime) Create(ctx context.Context, workload workload.Workload) (string, error) {

	labels := make(map[string]string)
	labels["orchestrator.workload.id"] = workload.ID
	labels["orchestrator.namespace"] = workload.Namespace
	labels["orchestrator.displayName"] = workload.DisplayName

	cfg := container.Config{
		Image:  workload.Image,
		Labels: labels,
	}

	hostCfg := container.HostConfig{
		NanoCPUs: workload.Resources.CPU * 1e9,
		Memory:   workload.Resources.Mem * 1024 * 1024,
	}

	result, err := d.apiClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:           &cfg,
		HostConfig:       &hostCfg,
		NetworkingConfig: nil,
		Platform:         nil,
		Name:             "merlin-" + workload.ID,
		Image:            workload.Image,
	})
	if err != nil {
		return "", err
	}

	log.Printf("Created container %s", result.ID)
	return result.ID, nil
}

// Start starts the container with the given containerID using Docker.
func (d *DockerRuntime) Start(ctx context.Context, containerID string) error {
	_, err := d.apiClient.ContainerStart(ctx, containerID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	log.Printf("Started container %s", containerID)
	return nil
}

// Stop stops the container with the given containerID using Docker.
func (d *DockerRuntime) Stop(ctx context.Context, containerID string) error {
	_, err := d.apiClient.ContainerStop(ctx, containerID, client.ContainerStopOptions{})
	if err != nil {
		return err
	}
	log.Printf("Stopped container %s", containerID)
	return nil
}

// Delete deletes the container with the given containerID using Docker.
func (d *DockerRuntime) Delete(ctx context.Context, containerID string) error {
	_, err := d.apiClient.ContainerRemove(ctx, containerID, client.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	log.Printf("Deleted container %s", containerID)
	return nil
}

// Inspect returns information about the container and workload associated with the given containerID.
func (d *DockerRuntime) Inspect(ctx context.Context, containerID string) (ContainerInfo, error) {
	result, err := d.apiClient.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return ContainerInfo{}, err
	}

	var info ContainerInfo

	info.ContainerID = result.Container.ID
	info.WorkloadID = result.Container.Config.Labels["orchestrator.workload.id"]

	info.ResourceLimits.CPU = result.Container.HostConfig.NanoCPUs / 1e9
	info.ResourceLimits.Mem = result.Container.HostConfig.Memory / 1024 / 1024

	info.ExitCode = result.Container.State.ExitCode
	info.Image = result.Container.Config.Image
	info.Status = ContainerStatus(result.Container.State.Status)
	info.Error = result.Container.State.Error

	return info, nil
}

// List provides a map of all Docker containers managed by the program with their containerInfo as the value.
func (d *DockerRuntime) List(ctx context.Context) (map[string]ContainerInfo, error) {
	list, err := d.apiClient.ContainerList(ctx, client.ContainerListOptions{All: true, Filters: make(client.Filters).Add("orchestrator.workload.id")})
	if err != nil {
		return nil, err
	}

	result := make(map[string]ContainerInfo)

	for _, ctr := range list.Items {
		info, err := d.Inspect(ctx, ctr.ID)
		if err != nil {
			return nil, err
		}

		result[ctr.ID] = info
	}

	return result, nil
}
