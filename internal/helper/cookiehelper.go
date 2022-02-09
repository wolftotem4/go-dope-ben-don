package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/http"

	"github.com/pkg/errors"
)

func GetFirstCookieValue(entries []*http.Cookie, name string) (string, bool) {
	for _, entry := range entries {
		if entry.Name == name {
			return entry.Value, true
		}
	}
	return "", false
}

func SerializeCookie(cookie *http.Cookie) (string, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(cookie); err != nil {
		return "", errors.WithStack(err)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func UnserializeCookie(raw string) (*http.Cookie, error) {
	b, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buffer := bytes.Buffer{}
	if _, err := buffer.Write(b); err != nil {
		return nil, errors.WithStack(err)
	}

	decoder := gob.NewDecoder(&buffer)

	var cookie http.Cookie
	if err := decoder.Decode(&cookie); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cookie, nil
}

func GetCookieNameSlice(cookies []*http.Cookie) []string {
	var names = make([]string, len(cookies))
	for i, cookie := range cookies {
		names[i] = cookie.Name
	}
	return names
}
