package humanize

import (
	"context"
	"math/rand"
	"time"

	"github.com/playwright-community/playwright-go"
	"golang.org/x/time/rate"
)

type Behavior struct {
	limiter  *rate.Limiter
	jitter   time.Duration
	patterns []Pattern
}

type Pattern struct {
	Action      string
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
				Action:      "scroll",
				MinDelay:    time.Second * 1,
				MaxDelay:    time.Second * 3,
				Probability: 0.7,
			},
			{
				Action:      "mousemove",
				MinDelay:    time.Millisecond * 500,
				MaxDelay:    time.Second * 2,
				Probability: 0.5,
			},
			{
				Action:      "pause",
				MinDelay:    time.Second * 3,
				MaxDelay:    time.Second * 8,
				Probability: 0.3,
			},
			{
				Action:      "highlight",
				MinDelay:    time.Millisecond * 200,
				MaxDelay:    time.Second * 1,
				Probability: 0.2,
			},
			{
				Action:      "viewport_scroll",
				MinDelay:    time.Second * 2,
				MaxDelay:    time.Second * 5,
				Probability: 0.4,
			},
		},
	}
}

func (b *Behavior) Wait(ctx context.Context) error {
	err := b.limiter.Wait(ctx)
	if err != nil {
		return err
	}

	jitter := time.Duration(rand.Int63n(int64(b.jitter)))
	time.Sleep(jitter)

	return nil
}

func (b *Behavior) SimulateHumanBehavior(page playwright.Page) error {
	for _, pattern := range b.patterns {
		if rand.Float64() > pattern.Probability {
			continue
		}

		delay := pattern.MinDelay + time.Duration(rand.Float64()*float64(pattern.MaxDelay-pattern.MinDelay))
		time.Sleep(delay)

		switch pattern.Action {
		case "scroll":
			_, err := page.Evaluate(`window.scrollBy({
                top: Math.floor(Math.random() * 100),
                behavior: 'smooth'
            })`)
			if err != nil {
				return err
			}
		case "mousemove":
			err := page.Mouse().Move(
				rand.Float64()*1000,
				rand.Float64()*800,
				playwright.MouseMoveOptions{Steps: playwright.Int(rand.Intn(5) + 3)},
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
