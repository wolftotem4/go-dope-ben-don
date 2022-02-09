package form

import (
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

type LoginForm struct {
	Username string `url:"username"`
	Password string `url:"password"`
	Answer   string `url:"result"`
	Submit   string `url:"submit"`
}

func (form *LoginForm) ToUrlValuesObj() (url.Values, error) {
	values, err := query.Values(form)
	if err != nil {
		return url.Values{}, errors.WithStack(err)
	}

	values.Set("signInPanel_signInForm:hf:0", "")
	values.Set("rememberMeRow:rememberMe", "on")
	values.Set("submit", "登入")
	return values, nil
}
