package action

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/config"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
	apperrors "github.com/wolftotem4/go-dope-ben-don/internal/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/form"
	"github.com/wolftotem4/go-dope-ben-don/internal/parser"
)

type LoginAction struct {
	Client client.Client
	Config *config.App
}

type LoggedInData struct {
	Interface int
}

func PerformLoginIfRequired(ctx context.Context, client client.Client, config *config.App, err error) (*LoggedInData, error) {
	var loginRequiredErr *apperrors.LoginRequiredError
	if errors.As(err, &loginRequiredErr) {
		action := &LoginAction{Client: client, Config: config}
		return action.PerformLogin(ctx, loginRequiredErr.Form)
	}
	return nil, err
}

func (action *LoginAction) GetFormData(ctx context.Context) (*parser.LoginForm, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := action.Client.Call(ctx, "/do/login", client.Webpage)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	data, err := action.ParseResponse(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (action *LoginAction) request(data *parser.LoginForm) (*http.Request, error) {
	answer, err := parser.AnswerQuestion(data.Question)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	f := &form.LoginForm{
		Username: action.Config.Account,
		Password: action.Config.Password,
		Answer:   answer,
	}

	params, err := f.ToUrlValuesObj()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return action.Client.NewPostRequest(data.Action, params)
}

func (action *LoginAction) PerformLogin(ctx context.Context, data *parser.LoginForm) (*LoggedInData, error) {
	req, err := action.request(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := action.Client.Send(ctx, req, client.Webpage)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	if err := handleFeedbackErrMsg(body); err != nil {
		return nil, err
	}

	interfaceValue, err := parser.ExtractInterface(body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &LoggedInData{
		Interface: interfaceValue,
	}, nil
}

func (action *LoginAction) ParseResponse(res *http.Response) (*parser.LoginForm, error) {
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return parser.ParseLoginHTML(html)
}
