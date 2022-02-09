package api

import "errors"

var ErrRefreshPageRequired = errors.New("refresh page is required due to API expiration")

func handleErrorMsg(msg string) error {
	switch msg {
	case "閒置過久，網頁已過期":
		return ErrRefreshPageRequired
	default:
		return errors.New(msg)
	}
}
