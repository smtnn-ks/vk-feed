package imagechecker

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// TODO: make richer errors
var ErrBadImage error = errors.New("something wrong with image")

type IC struct{}

func (ic IC) Check(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return ErrBadImage
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrBadImage
	}
	contentType := res.Header.Get("content-type")
	if contentType == "" {
		return ErrBadImage
	}
	if strings.Split(contentType, "/")[0] != "image" {
		return ErrBadImage
	}
	if res.ContentLength >= 5*10e6 {
		return ErrBadImage
	}
	return nil
}
