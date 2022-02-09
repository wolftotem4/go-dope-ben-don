package client

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/helper"
)

type MockClient struct {
	res  []*http.Response
	Base *url.URL
}

func NewMockClient() *MockClient {
	return &MockClient{
		Base: rootUrl,
	}
}

func (client *MockClient) AppendResponse(res *http.Response) {
	client.res = append(client.res, res)
}

func (client *MockClient) MockRemaining() int {
	return len(client.res)
}

func (client *MockClient) Load(ctx context.Context) error {
	panic("unimplemented")
}

func (client *MockClient) Send(ctx context.Context, req *http.Request, targetType int) (*http.Response, error) {
	if len(client.res) == 0 {
		return nil, errors.New("mocking failed, there is no response set")
	}

	res := client.res[0]
	client.res = client.res[1:]

	return res, nil
}

func (client *MockClient) Call(ctx context.Context, api string, targetType int) (*http.Response, error) {
	if len(client.res) == 0 {
		return nil, errors.New("mocking failed, there is no response set")
	}

	res := client.res[0]
	client.res = client.res[1:]

	return res, nil
}

func (client *MockClient) NewPostRequest(api string, params url.Values) (*http.Request, error) {
	u, err := client.Endpoint(api)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return helper.BuildPostRequest(u.String(), params)
}

func (client *MockClient) Endpoint(api string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimLeft(api, "/"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client.Base.ResolveReference(u), nil
}
