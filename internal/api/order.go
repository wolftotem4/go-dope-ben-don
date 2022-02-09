package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	typesjson "github.com/wolftotem4/go-dope-ben-don/internal/types/json"
)

func CallProductGroupForBuyer(ctx context.Context, orderId string, c client.Client) (*typesjson.ProductGroupForBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	path := fmt.Sprintf("/mvc/api/order/%s/product-group-for-buyer?expand=false", url.PathEscape(orderId))
	res, err := c.Call(ctx, path, client.Api)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res.Body.Close()

	var data typesjson.ProductGroupForBuyer
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.WithStack(err)
	}

	return &data, nil
}
