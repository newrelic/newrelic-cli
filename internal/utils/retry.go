package utils

import (
	"context"
	"time"
)

type Retry struct {
	MaxRetries         int
	retryDelayMilliSec int
	RetryFunc          func() error
}

func NewRetry(maxRetries int, retryDelayMilliSec int, retryFunc func() error) *Retry {
	return &Retry{
		MaxRetries:         maxRetries,
		retryDelayMilliSec: retryDelayMilliSec,
		RetryFunc:          retryFunc,
	}
}

func (r *Retry) ExecWithRetries(ctx context.Context) error {
	tries := 0
	success := false
	for !success {
		tries++
		if err := r.RetryFunc(); err != nil {
			if tries == r.MaxRetries {
				return err
			}

			w := make(chan struct{}, 1)
			go func() {
				time.Sleep(time.Duration(r.retryDelayMilliSec) * time.Millisecond)
				w <- struct{}{}
			}()

			select {
			case <-ctx.Done():
				return context.Canceled
			case <-w:
			}
		} else {
			success = true
		}
	}

	return nil
}
