package action

import (
	"context"
	"testing"

	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	apperrors "github.com/wolftotem4/go-dope-ben-don/internal/errors"
)

func TestDashboardLoginRequiredAction(t *testing.T) {
	mock := client.NewMockClient()

	res, err := createHtmlFileResponse("GET", "https://dinbendon.net/do/login", "testdata/login.htm")
	if err != nil {
		t.Fatal(err)
	}

	mock.AppendResponse(res)

	act := &DashboardAction{Client: mock}
	_, err = act.Update(context.Background())
	if _, isLoginRequiredErr := err.(*apperrors.LoginRequiredError); !isLoginRequiredErr {
		t.Fatal("expected LoginRequiredError")
	}
}
