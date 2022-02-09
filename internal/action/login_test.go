package action

import (
	"context"
	"testing"

	"github.com/wolftotem4/go-dope-ben-don/config"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	"github.com/wolftotem4/go-dope-ben-don/internal/parser"
)

func TestGetFormData(t *testing.T) {
	html, err := embedFS.ReadFile("testdata/login.htm")
	if err != nil {
		t.Fatal(err)
	}

	mock := client.NewMockClient()

	mock.AppendResponse(createResponse("GET", "https://dinbendon.net/do/login", html))

	act := &LoginAction{Client: mock}
	form, err := act.GetFormData(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if form.Action != "/do/?wicket:interface=:1:signInPanel:signInForm::IFormSubmitListener" {
		t.Errorf("unexpected action: %s", form.Action)
	}

	if form.Question != "99＋0＝" {
		t.Errorf("unexpected question: %s", form.Question)
	}
}

func TestLoginErrMsg(t *testing.T) {
	html, err := embedFS.ReadFile("testdata/login_errmsg.htm")
	if err != nil {
		t.Fatal(err)
	}

	mock := client.NewMockClient()

	mock.AppendResponse(createResponse("POST", "https://dinbendon.net/do/?wicket:interface=:1:signInPanel:signInForm::IFormSubmitListener", html))

	act := &LoginAction{Client: mock, Config: &config.App{
		Account:  "test",
		Password: "test",
		Name:     "test",
	}}
	_, err = act.PerformLogin(context.Background(), &parser.LoginForm{
		Logged:   false,
		Action:   "https://dinbendon.net/do/?wicket:interface=:1:signInPanel:signInForm::IFormSubmitListener",
		Question: "36+1=",
	})

	if err == nil {
		t.Fatal("expected error")
	} else if err.Error() != "Wrong password" {
		t.Errorf("unexpected error msg: %s", err.Error())
	}
}
