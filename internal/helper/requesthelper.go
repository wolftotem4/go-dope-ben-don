package helper

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func BuildPostRequest(url string, params url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, io.NopCloser(strings.NewReader(params.Encode())))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
