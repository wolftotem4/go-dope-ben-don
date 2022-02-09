package parser

import "testing"

func TestExtractInterface(t *testing.T) {
	body, err := embedFS.ReadFile("testdata/login_logged.htm")
	if err != nil {
		t.Fatal(err)
	}

	value, err := ExtractInterface(body)
	if err != nil {
		t.Fatal(err)
	}

	if value != 4 {
		t.Errorf("unexpected interface: %d", value)
	}
}
