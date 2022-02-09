package action

import (
	"bytes"
	"embed"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

//go:embed testdata/*
var embedFS embed.FS

func createResponse(method string, uri string, body []byte) *http.Response {
	u, _ := url.Parse(uri)
	return &http.Response{
		StatusCode: 200,
		Request: &http.Request{
			Method: method,
			URL:    u,
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	}
}

func createHtmlFileResponse(method string, uri string, file string) (*http.Response, error) {
	html, err := embedFS.ReadFile(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return createResponse(method, uri, html), nil
}
