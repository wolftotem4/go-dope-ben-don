package client

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/helper"
)

var Webpage = 1
var Api = 2
var WebpageApi = 3 // without "x-dbd-xsrf"

type Client interface {
	Load(ctx context.Context) error
	Call(ctx context.Context, api string, targetType int) (*http.Response, error)
	Send(ctx context.Context, req *http.Request, targetType int) (*http.Response, error)
	Endpoint(api string) (*url.URL, error)
	NewPostRequest(api string, params url.Values) (*http.Request, error)
}

type HttpClient struct {
	Http   *http.Client
	Base   *url.URL
	Status *Status
	Store  Store
}

var rootUrl, _ = url.Parse("https://dinbendon.net/")

func NewClient(store Store) (*HttpClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	http := &http.Client{Jar: jar}

	return &HttpClient{
		Http:   http,
		Base:   rootUrl,
		Status: NewStatus(),
		Store:  store,
	}, nil
}

func (client *HttpClient) Load(ctx context.Context) error {
	if err := client.Store.RestoreCookies(ctx); err != nil {
		return errors.WithStack(err)
	}

	client.Http.Jar.SetCookies(rootUrl, client.Store.GetCookies())

	return nil
}

func (client *HttpClient) updateCookies(res *http.Response) {
	client.Store.SetCookies(client.Http.Jar.Cookies(res.Request.URL))
}

func (client *HttpClient) Send(ctx context.Context, req *http.Request, targetType int) (*http.Response, error) {
	switch targetType {
	case Api:
		token, _ := helper.GetFirstCookieValue(client.Http.Jar.Cookies(req.URL), "DBD-XSRF")
		req.Header.Set("x-dbd-xsrf", token)
	}

	res, err := client.Http.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client.updateCookies(res)

	return res, err
}

func (client *HttpClient) Call(ctx context.Context, api string, targetType int) (*http.Response, error) {
	u, err := client.Endpoint(api)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client.Send(ctx, req, targetType)
}

func (client *HttpClient) NewPostRequest(api string, params url.Values) (*http.Request, error) {
	u, err := client.Endpoint(api)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return helper.BuildPostRequest(u.String(), params)
}

func (client *HttpClient) Endpoint(api string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimLeft(api, "/"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client.Base.ResolveReference(u), nil
}
