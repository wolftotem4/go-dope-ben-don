package typesjson

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalProductGroupForBuyer(t *testing.T) {
	j, err := embedFS.ReadFile("testdata/product-group-for-buyer.json")
	if err != nil {
		t.Fatal(err)
	}

	var data ProductGroupForBuyer

	if err := json.Unmarshal(j, &data); err != nil {
		t.Fatal(err)
	}

	if data.Data.Rows[0].Items[0].MergedName != "ThisIsAName" {
		t.Fatalf("unexpected value: %s", data.Data.Rows[0].Items[0].MergedName)
	}
}
