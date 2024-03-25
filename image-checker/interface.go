package imagechecker

import "context"

type ImageChecker interface {
	Check(ctx context.Context, url string) error
}
