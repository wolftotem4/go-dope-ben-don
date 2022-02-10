package gui

import (
	"log"
	"os"

	"github.com/asticode/go-astilectron"
	"github.com/pkg/errors"
)

type GUI struct {
	astilectron *astilectron.Astilectron
	tray        *astilectron.Tray
	dir         string

	Quit     chan bool
	EventErr chan *EventErr
}

func New(dir string, appName string) (*GUI, error) {
	a, err := astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:           appName,
		BaseDirectoryPath: dir,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &GUI{
		astilectron: a,
		dir:         dir,
		Quit:        make(chan bool),
		EventErr:    make(chan *EventErr),
	}, nil
}

func (gui *GUI) HandleSignals() {
	gui.astilectron.HandleSignals()
}

func (gui *GUI) Build() error {
	gui.astilectron.Start()

	if err := gui.buildTray(); err != nil {
		errors.WithStack(err)
	}

	if err := gui.buildMenu(); err != nil {
		errors.WithStack(err)
	}

	return nil
}

func (gui *GUI) Wait() error {
	gui.astilectron.Wait()

	return nil
}

func (gui *GUI) Close() {
	gui.astilectron.Close()
}
