package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/helper"
	"github.com/wolftotem4/go-dope-ben-don/internal/log"
)

type ClientLogger struct {
	Base Client
	Dir  string
}

func WrapLogger(client Client, dir string) *ClientLogger {
	return &ClientLogger{
		Base: client,
		Dir:  dir,
	}
}

func (client *ClientLogger) Load(ctx context.Context) error {
	return client.Base.Load(ctx)
}

func (client *ClientLogger) Call(ctx context.Context, api string, targetType int) (*http.Response, error) {
	res, err := client.Base.Call(ctx, api, targetType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := client.logResponse(res); err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (client *ClientLogger) Send(ctx context.Context, req *http.Request, targetType int) (*http.Response, error) {
	res, err := client.Base.Send(ctx, req, targetType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := client.logResponse(res); err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (client *ClientLogger) Endpoint(api string) (*url.URL, error) {
	return client.Base.Endpoint(api)
}

func (client *ClientLogger) NewPostRequest(api string, params url.Values) (*http.Request, error) {
	return client.Base.NewPostRequest(api, params)
}

func (client *ClientLogger) logResponse(res *http.Response) error {
	body, err := helper.ReadResponseBody(res)
	if err != nil {
		return errors.WithStack(err)
	}

	log.SaveResponse(res, body, client.Dir)

	return nil
}
