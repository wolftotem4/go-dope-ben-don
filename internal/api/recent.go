package api

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	typesjson "github.com/wolftotem4/go-dope-ben-don/internal/types/json"
)

func CallRecent(ctx context.Context, c client.Client) (*typesjson.Recent, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := c.Call(ctx, "/mvc/api/order/recent", client.Api)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	var data typesjson.Recent
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.WithStack(err)
	}

	if data.Error != "" {
		return nil, handleErrorMsg(data.Error)
	}

	return &data, nil
}
