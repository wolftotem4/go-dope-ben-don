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

// prompting items list and recommanded update time
type PromptItemInfo struct {
	// prompting items
	Items []*typesjson.ProgressItem

	// suggesting checkpoint to fetch next latest data
	UpdateTimer *time.Timer

	Interface int
}

// fetch prompting items and recommanded update time
func GetPromptItems(ctx context.Context, app *app.App, interfaceValue int) (*PromptItemInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var (
		pendingItems []*typesjson.ProgressItem
		closest      = newClosetRefreshTime(app.Config.PriorTime, app.Config.UpdateInterval)
	)

	if interfaceValue == 0 {
		// Refresh session status by loading the main page
		act := &action.DashboardAction{Client: app.Client}
		data, err := act.Update(ctx)
		// We don't need to handle the error here

		// Re-login if needed
		logged, err := action.PerformLoginIfRequired(ctx, app.Client, app.Config, err)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if logged != nil {
			interfaceValue = logged.Interface
		} else {
			interfaceValue = data.Interface
		}
	}

	// Load all pending items
	{
		data, err := api.CallProgress(ctx, app.Client)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		// convert to slice of pointers
		for _, item := range data.Data {
			pendingItems = append(pendingItems, &item)
		}
	}

	// Filter out ordered items
	{
		var filterItems []*typesjson.ProgressItem
		for _, item := range pendingItems {
			data, err := api.CallProductGroupForBuyer(ctx, item.OrderHashId, app.Client)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if !isItemOrdered(data, app.Config.Name) {
				filterItems = append(filterItems, item)
			}
		}
		pendingItems = filterItems
	}

	// Filter prompt items
	{
		var filterItems []*typesjson.ProgressItem
		for _, item := range pendingItems {
			if item.IsExpiring(app.Config.PriorTime) {
				filterItems = append(filterItems, item)
			}

			closest.update(item.RemainSecondBeforeExpire)
		}
		pendingItems = filterItems
	}

	return &PromptItemInfo{
		Items:       pendingItems,
		UpdateTimer: closest.timer(),
		Interface:   interfaceValue,
	}, nil
}

// determained if a item has been ordered by matching name
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

type closetRefreshTime struct {
	Time      time.Duration
	PriorTime time.Duration
	Interval  time.Duration
}

func newClosetRefreshTime(priorTime time.Duration, interval time.Duration) *closetRefreshTime {
	return &closetRefreshTime{
		PriorTime: priorTime,
		Interval:  interval,
	}
}

func (timer *closetRefreshTime) isSet() bool {
	return timer.Time > 0
}

func (timer *closetRefreshTime) update(newTime time.Duration) {
	newTime -= timer.PriorTime
	if newTime > 0 && (!timer.isSet() || newTime < timer.Time) {
		timer.Time = newTime
	}
}

func (timer *closetRefreshTime) next() time.Duration {
	if timer.isSet() && timer.Time < timer.Interval {
		return timer.Time
	} else {
		return timer.Interval
	}
}

func (timer *closetRefreshTime) timer() *time.Timer {
	return time.NewTimer(timer.next())
}
