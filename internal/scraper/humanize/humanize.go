package humanize

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/time/rate"
)

type Behavior struct {
	limiter  *rate.Limiter
	jitter   time.Duration
	patterns []Pattern
}

type Pattern struct {
	MinDelay    time.Duration
	MaxDelay    time.Duration
	Probability float64
}

func NewBehavior(rps float64, burst int) *Behavior {
	return &Behavior{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
		jitter:  time.Millisecond * 100,
		patterns: []Pattern{
			{
				MinDelay:    time.Second * 1,
				MaxDelay:    time.Second * 3,
				Probability: 0.7,
			},
			{
				MinDelay:    time.Millisecond * 500,
				MaxDelay:    time.Second * 2,
				Probability: 0.5,
			},
		},
	}
}

func (b *Behavior) Wait(ctx context.Context) error {
	if err := b.limiter.Wait(ctx); err != nil {
		return err
	}

	// Add random jitter
	jitter := time.Duration(rand.Int63n(int64(b.jitter)))
	time.Sleep(jitter)

	// Apply pattern-based delays
	for _, pattern := range b.patterns {
		if rand.Float64() > pattern.Probability {
			continue
		}
		delay := pattern.MinDelay + time.Duration(rand.Float64()*float64(pattern.MaxDelay-pattern.MinDelay))
		time.Sleep(delay)
	}

	return nil
}
