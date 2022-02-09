package action

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	apperrors "github.com/wolftotem4/go-dope-ben-don/internal/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/parser"
)

type DashboardAction struct {
	Client client.Client
}

type DashboardData struct {
	Interface int
}

func (action *DashboardAction) Update(ctx context.Context) (*DashboardData, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := action.Client.Call(ctx, "/do/", client.Webpage)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	if strings.HasPrefix(res.Request.URL.EscapedPath(), "/do/login") {
		form, err := parser.ParseLoginHTML(body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return nil, apperrors.NewLoginRequiredError("login required", form)
	}

	interfaceValue, err := parser.ExtractInterface(body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &DashboardData{
		Interface: interfaceValue,
	}, nil
}
