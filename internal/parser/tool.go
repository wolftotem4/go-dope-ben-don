package parser

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func ExtractInterface(html []byte) (int, error) {
	regex := regexp.MustCompile(`(?i)\?wicket:interface=:(\d+):`)
	match := regex.FindSubmatch(html)

	if len(match) == 0 {
		return 0, errors.New("no interface found")
	}

	value, err := strconv.Atoi(string(match[1]))
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return value, nil
}
