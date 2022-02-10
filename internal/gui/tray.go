package gui

import (
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/pkg/errors"
)

func (gui *GUI) buildTray() error {
	tray := gui.astilectron.NewTray(&astilectron.TrayOptions{
		Image:   astikit.StrPtr(filepath.Join(gui.dir, "assets", "information.png")),
		Tooltip: astikit.StrPtr("Test"),
	})
	if err := tray.Create(); err != nil {
		return errors.WithStack(err)
	}

	gui.tray = tray

	return nil
}

func (gui *GUI) buildMenu() error {
	var m = gui.tray.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("離開"),
		},
	})

	mi, err := m.Item(0)
	if err != nil {
		return errors.WithStack(err)
	}

	mi.On(astilectron.EventNameMenuItemEventClicked, func(e astilectron.Event) (deleteListener bool) {
		gui.Quit <- true
		gui.tray.Destroy()
		gui.astilectron.Quit()
		return
	})

	return m.Create()
}
