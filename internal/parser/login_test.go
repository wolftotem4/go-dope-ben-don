package parser

import (
	"testing"

	"github.com/pkg/errors"
)

func TestLoginParser(t *testing.T) {
	html, err := embedFS.ReadFile("testdata/login.htm")
	if err != nil {
		t.Fatal(err)
	}

	form, err := ParseLoginHTML(html)
	if err != nil {
		t.Fatal(err)
	}

	if form.Logged {
		t.Error("unexpected value (Logged: true)")
	}

	if form.Action != "/do/?wicket:interface=:1:signInPanel:signInForm::IFormSubmitListener" {
		t.Errorf("unexpected action: %s", form.Action)
	}

	if form.Question != "99＋0＝" {
		t.Errorf("unexpected question: %s", form.Question)
	}
}

func TestLoggedLoginParser(t *testing.T) {
	files := []string{"testdata/login_logged.htm", "testdata/login_english.htm"}

	for _, file := range files {
		html, err := embedFS.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		form, err := ParseLoginHTML(html)
		if err != nil {
			t.Fatal(err)
		}

		if !form.Logged {
			t.Errorf("unexpected value (%s, Logged: false)", file)
		}
	}
}

func TestParseQuestion(t *testing.T) {
	data := []string{
		"99＋0＝",
		"30 加 5 等於",
		"36+1=",
	}

	for _, str := range data {
		q, err := ParseQuestion([]byte(str))
		if err != nil {
			t.Error(errors.WithMessagef(err, "test failed (%s)", str))
		} else if q != str {
			t.Errorf("unexpected question: %s", q)
		}
	}
}

func TestAnswerQuestion(t *testing.T) {
	data := []struct {
		Question string
		Anwser   string
	}{
		{Question: "99＋1＝", Anwser: "100"},
		{Question: "30 加 5 等於", Anwser: "35"},
		{Question: "36+1=", Anwser: "37"},
	}

	for _, item := range data {
		answer, err := AnswerQuestion(item.Question)
		if err != nil {
			t.Error(err)
		} else if answer != item.Anwser {
			t.Errorf("unexpected answer: %s", answer)
		}
	}
}
