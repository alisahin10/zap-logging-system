package logswitcher

import "sync/atomic"

// ActiveFlag represents a thread-safe boolean flag that can be shared across multiple goroutines and components.
type ActiveFlag struct {
	flag atomic.Bool // Atomic boolean for lock-free thread-safe operations
}

// NewActiveFlag creates a new ActiveFlag instance with the specified initial state.
func NewActiveFlag(initial bool) *ActiveFlag {
	a := &ActiveFlag{}
	a.flag.Store(initial) // Atomically set the initial value
	return a
}

// IsActive returns the current state of the flag in a thread-safe manner.
func (a *ActiveFlag) IsActive() bool {
	return a.flag.Load() // Atomically load the current value
}

// Toggle atomically flips the current state of the flag from true to false or from false to true.
func (a *ActiveFlag) Toggle() {
	a.flag.Store(!a.IsActive())
}
