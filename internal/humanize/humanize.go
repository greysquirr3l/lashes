package humanize

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"golang.org/x/time/rate"
)

// Pattern defines a behavior pattern with delay characteristics and probability
type Pattern struct {
	MinDelay    time.Duration
	MaxDelay    time.Duration
	Probability float64
}

type Action struct {
	Type      string  // scroll, mousemove, click, wait
	X         float64 // coordinate for mouse actions
	Y         float64 // coordinate for mouse actions
	Duration  time.Duration
	Scrolling struct {
		DeltaY   int
		Behavior string
		Segments int
		Interval time.Duration
	}
}

type Behavior struct {
	limiter  *rate.Limiter
	jitter   time.Duration
	patterns []Pattern
	viewport struct {
		width  int
		height int
	}
}

// NewBehavior creates a new humanized behavior simulator
func NewBehavior(rps float64, burst int) *Behavior {
	b := &Behavior{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
		jitter:  time.Millisecond * 100,
		patterns: []Pattern{
			{MinDelay: time.Second, MaxDelay: time.Second * 3, Probability: 0.7},
			{MinDelay: time.Millisecond * 500, MaxDelay: time.Second * 2, Probability: 0.5},
		},
	}
	b.viewport.width = 1920
	b.viewport.height = 1080
	return b
}

// SimulateAction generates a human-like action
func (b *Behavior) SimulateAction() *Action {
	actions := []string{"scroll", "mousemove", "click", "wait"}
	actionType := secureSelectFromSlice(actions)

	action := &Action{
		Type:     actionType,
		Duration: b.randomDuration(200*time.Millisecond, 2*time.Second),
	}

	switch actionType {
	case "scroll":
		action.Scrolling.DeltaY = secureRandInt(300) + 100
		action.Scrolling.Behavior = "smooth"
		action.Scrolling.Segments = secureRandInt(5) + 3
		action.Scrolling.Interval = time.Millisecond * time.Duration(secureRandInt(50)+30)
	case "mousemove":
		action.X = secureRandFloat() * float64(b.viewport.width)
		action.Y = secureRandFloat() * float64(b.viewport.height)
	}

	return action
}

// Wait blocks for an appropriate duration based on the behavior pattern
func (b *Behavior) Wait(ctx context.Context) error {
	if err := b.limiter.Wait(ctx); err != nil {
		return err
	}

	// Add secure random jitter
	jitter := time.Duration(secureRandInt(int(b.jitter)))
	time.Sleep(jitter)

	// Apply pattern-based delays
	for _, pattern := range b.patterns {
		if secureRandFloat() < pattern.Probability {
			delay := b.randomDuration(pattern.MinDelay, pattern.MaxDelay)
			time.Sleep(delay)
			break
		}
	}

	return nil
}

// Helper functions using crypto/rand
func secureRandInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}

func secureRandFloat() float64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return 0.0
	}
	return float64(n.Int64()) / 1000000.0
}

func secureSelectFromSlice(items []string) string {
	idx := secureRandInt(len(items))
	return items[idx]
}

func (b *Behavior) randomDuration(min, max time.Duration) time.Duration {
	delta := max - min
	return min + time.Duration(secureRandInt(int(delta)))
}
