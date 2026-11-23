// Package workload defines data types which are core to the whole of the program.
package workload

// DesiredState refers to the desired state of a workload set by the controller.
// It will nearly always be defined by only "Running" or "Stopped."
type DesiredState string

const (
	Running DesiredState = "running"
	Stopped DesiredState = "stopped"
)

// Resources refers to required or available CPU cores and memory depending on the given context.
type Resources struct {
	CPU int64
	Mem int64
}

// Workload refers to any given containerized application and its desired configuration.
type Workload struct {
	ID             string
	Image          string
	Desired        DesiredState
	Resources      Resources
	ConfigRevision uint64
}
