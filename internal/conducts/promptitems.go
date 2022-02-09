package conducts

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/action"
	"github.com/wolftotem4/go-dope-ben-don/internal/api"
	"github.com/wolftotem4/go-dope-ben-don/internal/app"
	typesjson "github.com/wolftotem4/go-dope-ben-don/internal/types/json"
)

func GetPromptItems(ctx context.Context, app *app.App, refreshSession bool) (*[]typesjson.ProgressItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var (
		pendingItems *[]typesjson.ProgressItem
	)

	if refreshSession {
		// Refresh session status by loading the main page
		act := &action.DashboardAction{Client: app.Client}
		_, err := act.Update(ctx)
		// We don't need to handle the error here

		// Re-login if needed
		_, err = action.PerformLoginIfRequired(ctx, app.Client, app.Config, err)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Load all pending items
	{
		data, err := api.CallProgress(ctx, app.Client)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		pendingItems = &data.Data
	}

	// Filter prompt items
	{
		var filterItems []typesjson.ProgressItem
		for _, item := range *pendingItems {
			if item.IsExpiring(app.Config.PriorTime) {
				filterItems = append(filterItems, item)
			}
		}
		pendingItems = &filterItems
	}

	// Filter out ordered items
	{
		var filterItems []typesjson.ProgressItem
		for _, item := range *pendingItems {
			data, err := api.CallProductGroupForBuyer(ctx, item.OrderHashId, app.Client)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if !isItemOrdered(data, app.Config.Name) {
				filterItems = append(filterItems, item)
			}
		}
		pendingItems = &filterItems
	}

	return pendingItems, nil
}

func isItemOrdered(pendingItem *typesjson.ProductGroupForBuyer, name string) bool {
	for _, row := range pendingItem.Data.Rows {
		for _, subitem := range row.Items {
			if subitem.MergedName == name {
				return true
			}
		}
	}
	return false
}
