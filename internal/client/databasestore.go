package client

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/helper"
)

type DatabaseStore struct {
	DB        *sqlx.DB
	SessionId string
	Cookies   []*http.Cookie
}

func NewDatabaseStore(db *sqlx.DB) *DatabaseStore {
	return &DatabaseStore{
		DB: db,
	}
}

func (store *DatabaseStore) GetCookies() []*http.Cookie {
	return store.Cookies
}

func (store *DatabaseStore) SetCookies(cookies []*http.Cookie) {
	store.Cookies = cookies
}

func (store *DatabaseStore) RestoreCookies(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := store.DB.QueryContext(ctx, "SELECT value FROM cookies")
	if err != nil {
		return errors.WithStack(err)
	}

	var cookies []*http.Cookie
	for rows.Next() {
		var serialized string
		if err := rows.Scan(&serialized); err != nil {
			return errors.WithStack(err)
		}

		cookie, err := helper.UnserializeCookie(serialized)
		if err != nil {
			return errors.WithStack(err)
		}

		cookies = append(cookies, cookie)
	}

	store.Cookies = cookies

	return nil
}

func (store *DatabaseStore) Save(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := store.saveCookies(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (store *DatabaseStore) saveCookies(ctx context.Context) error {
	for _, cookie := range store.Cookies {
		raw, err := helper.SerializeCookie(cookie)
		if err != nil {
			return errors.WithStack(err)
		}

		if err := store.performCookieUpdate(ctx, cookie.Name, raw); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := store.excludeOtherCookies(ctx, helper.GetCookieNameSlice(store.Cookies)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (store *DatabaseStore) performCookieUpdate(ctx context.Context, name string, value string) error {
	stmt, err := store.DB.PrepareContext(ctx, "INSERT INTO cookies (name, value) VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET value = excluded.value")
	if err != nil {
		return errors.WithStack(err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, name, value); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (store *DatabaseStore) excludeOtherCookies(ctx context.Context, names []string) error {
	if len(names) == 0 {
		return nil
	}

	query, args, err := sqlx.In("DELETE FROM cookies WHERE name NOT IN (?)", names)
	if err != nil {
		return errors.WithStack(err)
	}

	query = store.DB.Rebind(query)
	if _, err := store.DB.ExecContext(ctx, query, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil

}
