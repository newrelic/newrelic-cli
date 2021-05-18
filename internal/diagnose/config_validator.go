package diagnose

import "context"

type ConfigValidator interface {
	ValidateConfig(ctx context.Context) error
}
