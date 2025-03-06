package captcha

import (
	"context"
	"errors"
	"time"

	"github.com/playwright-community/playwright-go"
)

type SolverType string

const (
	BrowserSolver SolverType = "browser"
	IPRotation    SolverType = "ip_rotation"
)

var (
	ErrCaptchaNotFound = errors.New("captcha not found")
	ErrSolveTimeout    = errors.New("captcha solve timeout")
	ErrUnsupported     = errors.New("unsupported captcha type")
)

type Options struct {
	Type          SolverType
	Timeout       time.Duration
	MaxAttempts   int
	RotateIPAfter int
	Debug         bool
}

type Solver interface {
	Solve(ctx context.Context, page playwright.Page) error
	HandleReCaptcha(ctx context.Context, page playwright.Page) error
	HandleHCaptcha(ctx context.Context, page playwright.Page) error
}

type solver struct {
	opts     Options
	attempts int
	lastIP   string
}

func NewSolver(opts Options) Solver {
	if opts.MaxAttempts == 0 {
		opts.MaxAttempts = 3
	}
	if opts.RotateIPAfter == 0 {
		opts.RotateIPAfter = 5
	}
	return &solver{opts: opts}
}

func (s *solver) Solve(ctx context.Context, page playwright.Page) error {
	// Try reCAPTCHA first
	if err := s.HandleReCaptcha(ctx, page); err == nil {
		return nil
	}

	// Try hCaptcha
	if err := s.HandleHCaptcha(ctx, page); err == nil {
		return nil
	}

	return ErrCaptchaNotFound
}

func (s *solver) HandleReCaptcha(ctx context.Context, page playwright.Page) error {
	// Try to find and click any "I'm not a robot" checkbox first
	if checkbox, err := page.QuerySelector(".recaptcha-checkbox"); err == nil && checkbox != nil {
		if err := checkbox.Click(); err == nil {
			// Wait to see if checkbox solving worked
			time.Sleep(2 * time.Second)
			return nil
		}
	}

	return ErrCaptchaNotFound
}

func (s *solver) HandleHCaptcha(ctx context.Context, page playwright.Page) error {
	// Similar simplified approach for hCaptcha
	if checkbox, err := page.QuerySelector(".checkbox"); err == nil && checkbox != nil {
		if err := checkbox.Click(); err == nil {
			time.Sleep(2 * time.Second)
			return nil
		}
	}

	return ErrCaptchaNotFound
}
