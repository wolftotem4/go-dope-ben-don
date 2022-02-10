package gui

import (
	"github.com/asticode/go-astilectron"
)

type EventErr struct {
	cause error
	Event astilectron.Event
}

func NewEventErr(event astilectron.Event, cause error) *EventErr {
	return &EventErr{
		cause: cause,
		Event: event,
	}
}

func (err EventErr) Error() string {
	return err.cause.Error()
}

func (err *EventErr) Cause() error {
	return err.cause
}

func (err *EventErr) Unwrap() error {
	return err.cause
}
