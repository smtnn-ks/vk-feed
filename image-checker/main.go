package imagechecker

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

var ErrUrlUnavailable error = errors.New("image url unavailable")
var ErrNotImage error = errors.New("image url leads to non-image content type")
var ErrImageTooBig error = errors.New("image too big")

type IC struct{}

func (ic IC) Check(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return ErrUrlUnavailable
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrUrlUnavailable
	}
	contentType := res.Header.Get("content-type")
	if contentType == "" {
		return ErrNotImage
	}
	if strings.Split(contentType, "/")[0] != "image" {
		return ErrNotImage
	}
	if res.ContentLength >= 5*10e6 {
		return ErrImageTooBig
	}
	return nil
}
