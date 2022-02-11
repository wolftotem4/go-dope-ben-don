package typesjson

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalProgress(t *testing.T) {
	var data Progress

	j, err := embedFS.ReadFile("testdata/progress.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(j, &data); err != nil {
		t.Fatal(err)
	}

	if len(data.Data) == 0 {
		t.Fatal("failed to unmarshal JSON")
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		t.Fatal(err)
	}

	if !time.Date(2022, 2, 8, 10, 0, 0, 0, loc).Equal(data.Data[0].ExpireDate.Time) {
		t.Fatalf("unexpected date: %s", data.Data[0].ExpireDate.Time.String())
	}
}

func TestIsExpiring(t *testing.T) {
	var item = &ProgressItem{RemainSecondBeforeExpire: 5 * 60}

	if !item.IsExpiring(6 * time.Minute) {
		t.Error("expect positive value")
	}
	if !item.IsExpiring(5 * time.Minute) {
		t.Error("expect positive value")
	}
	if item.IsExpiring(4 * time.Minute) {
		t.Error("expect negative value")
	}
	if item.IsExpiring(0) {
		t.Error("expect negative value")
	}
}
