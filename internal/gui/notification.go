package gui

import (
	"fmt"
	"path/filepath"

	"github.com/asticode/go-astilectron"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

type PromptItemInfo struct {
	Name string
	Url  string
}

func (gui *GUI) Notify(item PromptItemInfo) error {
	var n = gui.astilectron.NewNotification(&astilectron.NotificationOptions{
		Body:  fmt.Sprintf("[%s] 即將過期", item.Name),
		Icon:  filepath.Join(gui.dir, "assets", "information.png"),
		Title: "夠多便當",
	})

	n.On(astilectron.EventNameNotificationEventClicked, func(e astilectron.Event) (deleteListener bool) {
		if err := browser.OpenURL(item.Url); err != nil {
			gui.EventErr <- NewEventErr(e, errors.WithStack(err))
		}
		return
	})

	if err := n.Create(); err != nil {
		return errors.WithStack(err)
	}

	if err := n.Show(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
