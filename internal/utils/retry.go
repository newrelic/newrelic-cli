package utils

import (
	"context"
	"time"

	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

type RetryContext struct {
	RetryCount int
	Errors     []error
	Success    bool
	Canceled   bool
}

func (c *RetryContext) MostRecentError() error {
	if len(c.Errors) > 0 {
		return c.Errors[len(c.Errors)-1]
	}

	return nil
}

type Retry struct {
	MaxRetries   int
	retryDelayMs int
	RetryFunc    func() error
}

func NewRetry(maxRetries int, retryDelayMs int, retryFunc func() error) *Retry {
	return &Retry{
		MaxRetries:   maxRetries,
		retryDelayMs: retryDelayMs,
		RetryFunc:    retryFunc,
	}
}

func (r *Retry) ExecWithRetries(ctx context.Context) *RetryContext {
	retryCtx := RetryContext{}
	for !retryCtx.Success {
		retryCtx.RetryCount++
		if err := r.RetryFunc(); err != nil {
			if _, ok := err.(*nrErrors.PaymentRequiredError); ok {
				retryCtx.Success = false
				retryCtx.Errors = append(retryCtx.Errors, err)
				return &retryCtx
			}

			retryCtx.Errors = append(retryCtx.Errors, err)

			if retryCtx.RetryCount == r.MaxRetries {
				retryCtx.Success = false
				return &retryCtx
			}

			w := make(chan struct{}, 1)
			go func() {
				time.Sleep(time.Duration(r.retryDelayMs) * time.Millisecond)
				w <- struct{}{}
			}()

			select {
			case <-ctx.Done():
				retryCtx.Canceled = true
				retryCtx.Errors = append(retryCtx.Errors, context.Canceled)
				return &retryCtx
			case <-w:
			}
		} else {
			retryCtx.Success = true
		}
	}

	return &retryCtx
}
