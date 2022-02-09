package typesjson

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type Progress struct {
	Data  []ProgressItem `json:"data"`
	Error string         `json:"error,omitempty"`
}

type ProgressItem struct {
	Announcement             string    `json:"announcement"`
	ExpireDate               time.Time `json:"expireDate"`
	GroupId                  int       `json:"groupId"`
	InProgress               bool      `json:"inProgress"`
	MaxQty                   int       `json:"maxQty"`
	MaxTotalCost             int       `json:"maxTotalCost"`
	OrderHashId              string    `json:"orderHashId"`
	Originator               string    `json:"originator"`
	PasswordLocked           bool      `json:"passwordLocked"`
	RemainSecondBeforeExpire int       `json:"remainSecondBeforeExpire"`
	ShopName                 string    `json:"shopName"`
	Size                     int       `json:"size"`
	Total                    int       `json:"total"`
	Unlockable               bool      `json:"unlockable"`
}

func (item *ProgressItem) IsExpiring(priorTime time.Duration) bool {
	return time.Duration(item.RemainSecondBeforeExpire)*time.Second <= priorTime
}

func (item *ProgressItem) GetPath() string {
	return fmt.Sprintf("/do/order?id=%s", url.QueryEscape(item.OrderHashId))
}

func (item *ProgressItem) UnmarshalJSON(data []byte) error {
	type Alias ProgressItem

	aux := struct {
		ExpireDate int64 `json:"expireDate"`
		*Alias
	}{
		Alias: (*Alias)(item),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return errors.WithStack(err)
	}

	item.ExpireDate = time.UnixMilli(aux.ExpireDate)

	return nil
}
