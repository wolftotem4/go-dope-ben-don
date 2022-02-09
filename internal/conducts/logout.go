package conducts

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/action"
	"github.com/wolftotem4/go-dope-ben-don/internal/app"
)

func Logout(ctx context.Context, app *app.App) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var (
		interfaceValue int
	)

	// Refresh session status
	{
		// Refresh session status by loading the main page
		act := &action.DashboardAction{Client: app.Client}
		data, err := act.Update(ctx)
		if err != nil {
			return errors.WithStack(err)
		}

		// Get the interface value from the page
		interfaceValue = data.Interface
	}

	// Perform loggout
	{
		act := &action.LogoutAction{Client: app.Client}
		if err := act.ForceLogout(ctx, interfaceValue); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
