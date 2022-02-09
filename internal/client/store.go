package client

import (
	"context"
	"net/http"
)

type Store interface {
	RestoreCookies(ctx context.Context) error
	GetCookies() []*http.Cookie
	SetCookies(cookies []*http.Cookie)
	Save(ctx context.Context) error
}
