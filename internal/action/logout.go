package action

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	"github.com/wolftotem4/go-dope-ben-don/internal/parser"
)

type LogoutAction struct {
	Client client.Client
}

func (action *LogoutAction) PerformLogout(ctx context.Context, interfaceValue int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logoutUrl := fmt.Sprintf("/do/?wicket:interface=:%d:logout::ILinkListener", interfaceValue)

	res, err := action.Client.Call(ctx, logoutUrl, client.Webpage)
	if err != nil {
		return errors.WithStack(err)
	}

	if !strings.HasPrefix(res.Request.URL.EscapedPath(), "/do/login") {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.WithStack(err)
		}
		defer res.Body.Close()

		matched, err := regexp.Match(`(?i)Page Expired|網頁已過期`, body)
		if err != nil {
			return errors.WithStack(err)
		}

		if matched {
			return ErrPageExpired
		} else {
			return errors.New("failed to logout")
		}
	}

	return nil
}

func (action *LogoutAction) ForceLogout(ctx context.Context, interfaceValue int) error {
	for retry := 1; retry >= 0; retry-- {
		err := action.PerformLogout(ctx, interfaceValue)
		if err == nil {
			break
		} else if retry == 0 || err != ErrPageExpired {
			return err
		}

		// refresh homepage in order to renew expired Logout page
		res, err := action.Client.Call(ctx, "/do/", client.Webpage)
		if err != nil {
			return errors.WithStack(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.WithStack(err)
		}

		interfaceValue, err = parser.ExtractInterface(body)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
