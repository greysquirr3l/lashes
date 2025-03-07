// Package humanize provides a compatibility layer that re-exports
// functionality from the main humanize package
package humanize

import (
	"github.com/greysquirr3l/lashes/internal/humanize"
)

// Re-export needed types
type (
	Pattern  = humanize.Pattern
	Behavior = humanize.Behavior
)

// NewBehavior creates a new humanized behavior simulator
// This is a compatibility wrapper around the main humanize package
func NewBehavior(rps float64, burst int) *Behavior {
	return humanize.NewBehavior(rps, burst)
}
