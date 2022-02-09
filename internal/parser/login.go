package parser

import (
	"html"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type LoginForm struct {
	Logged   bool
	Action   string
	Question string
}

func ParseLoginHTML(data []byte) (*LoginForm, error) {
	var logged bool
	var action string

	logged, _ = regexp.Match(`Logout|登出`, data)

	{
		regex := regexp.MustCompile(`(?i)<form action="([^"]+)"`)
		match := regex.FindSubmatch(data)
		if len(match) == 0 {
			return nil, errors.New("no form action has been found")
		}
		action = html.UnescapeString(string(match[1]))
	}

	question, err := ParseQuestion(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &LoginForm{
		Logged:   logged,
		Action:   action,
		Question: question,
	}, nil
}

func ParseQuestion(data []byte) (string, error) {
	regex := regexp.MustCompile(`\d+(?:\+|＋| 加 )\d+(?:=|＝| 等於)`)
	match := regex.Find(data)
	if len(match) == 0 {
		return "", errors.New("no question has been found")
	}
	return string(match), nil
}

func AnswerQuestion(question string) (string, error) {
	regex := regexp.MustCompile(`(\d+)(\+|＋| 加 )(\d+)(?:=|＝| 等於)`)
	match := regex.FindStringSubmatch(question)
	if len(match) == 0 {
		return "", errors.Errorf("failed to parse question: %s", question)
	}

	left, err := strconv.Atoi(match[1])
	if err != nil {
		return "", errors.Wrap(err, "failed to parse question")
	}

	right, err := strconv.Atoi(match[3])
	if err != nil {
		return "", errors.Wrap(err, "failed to parse question")
	}

	operator := match[2]

	switch operator {
	case "+", "＋", " 加 ":
		return strconv.Itoa(left + right), nil
	default:
		return "", errors.Errorf("unreconized operator: %s", operator)
	}
}
