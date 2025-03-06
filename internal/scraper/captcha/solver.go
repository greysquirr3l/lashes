package captcha

import (
	"context"
	"errors"
	"time"
)

var (
    ErrCaptchaNotFound = errors.New("captcha not found")
    ErrSolveTimeout    = errors.New("captcha solve timeout")
    ErrUnsupported     = errors.New("unsupported captcha type")
)

// Config defines solver configuration
type Config struct {
    Timeout     time.Duration
    MaxAttempts int
    Debug       bool
}

// Result contains captcha solving result
type Result struct {
    Token     string
    Timestamp time.Time
    Duration  time.Duration
}

// Solver defines the interface for captcha solving
type Solver interface {
    // Solve attempts to solve a captcha on the given page
    Solve(ctx context.Context, pageURL string) (*Result, error)
}

// Note: Implementation should be provided by the end user
// This keeps the core library dependency-free
