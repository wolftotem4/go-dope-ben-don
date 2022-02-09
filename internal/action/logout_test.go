package action

import (
	"context"
	"testing"

	"github.com/wolftotem4/go-dope-ben-don/internal/client"
)

func TestLogoutPageExpired(t *testing.T) {
	files := []string{"testdata/logout_page_expired.htm", "testdata/logout_page_expired_2.htm"}

	for _, file := range files {
		html, err := embedFS.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		mock := client.NewMockClient()

		mock.AppendResponse(createResponse("GET", "https://dinbendon.net/do/?wicket:interface=:2:logout::ILinkListener", html))

		act := &LogoutAction{Client: mock}
		if err = act.PerformLogout(context.Background(), 2); err != ErrPageExpired {
			t.Error("expected page expiry")
		}
	}
}

func TestPerformLogout(t *testing.T) {
	html, err := embedFS.ReadFile("testdata/login.htm")
	if err != nil {
		t.Fatal(err)
	}

	mock := client.NewMockClient()

	mock.AppendResponse(createResponse("GET", "https://dinbendon.net/do/login", html))

	act := &LogoutAction{Client: mock}
	if err = act.PerformLogout(context.Background(), 13); err != nil {
		t.Error(err)
	}
}

func TestForceLogout(t *testing.T) {
	mock := client.NewMockClient()

	data := []struct {
		File string
		Url  string
	}{
		{File: "testdata/logout_page_expired.htm", Url: "https://dinbendon.net/do/?wicket:interface=:2:logout::ILinkListener"},
		{File: "testdata/login.htm", Url: "https://dinbendon.net/do/login"},
		{File: "testdata/login.htm", Url: "https://dinbendon.net/do/login"},
	}
	for _, item := range data {
		html, err := embedFS.ReadFile(item.File)
		if err != nil {
			t.Fatal(err)
		}

		mock.AppendResponse(createResponse("GET", item.Url, html))
	}

	act := &LogoutAction{Client: mock}
	if err := act.ForceLogout(context.Background(), 2); err != nil {
		t.Error(err)
	}

	if mock.MockRemaining() != 0 {
		t.Error("unexpected Login flow")
	}
}
