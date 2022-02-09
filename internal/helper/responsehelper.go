package helper

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func ReadResponseBody(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	defer res.Body.Close()

	// replace body to make sure it is readable
	res.Body = io.NopCloser(bytes.NewReader(body))

	return body, nil
}
