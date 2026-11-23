// Package runtime handles a runtime interface and its implementation(s) for containerized applications.
package runtime

import "github.com/will-petersmeyer/workload-orchestrator/internal/workload"

type ContainerRuntime interface {
	Create(workload workload.Workload) (string, error)
	Start(containerID string) error
	Stop(containerID string) error
	Delete(containerID string) error
	Inspect(containerID string) (ContainerInfo, error)
	List() ([]string, error)
}

type ContainerStatus string

const (
	Running    ContainerStatus = "running"
	Created    ContainerStatus = "created"
	Paused     ContainerStatus = "paused"
	Dead       ContainerStatus = "dead"
	Exited     ContainerStatus = "exited"
	Removing   ContainerStatus = "removing"
	Restarting ContainerStatus = "restarting"
)

type ContainerInfo struct {
	ContainerID    string
	WorkloadID     string
	Status         ContainerStatus
	Image          string
	ResourceLimits workload.Resources
	ExitCode       int
	Error          string
}
