package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type App struct {
	// Account name
	Account string

	// Password
	Password string

	// Name for determining if an item is ordered
	Name string

	// How long before to inform user there is prompt items
	PriorTime time.Duration

	// How often to fetch latest data
	UpdateInterval time.Duration

	Debug bool
}

const (
	EnvFileName = ".env"
)

func LoadConfig(dir string) (*App, error) {
	var (
		account   string
		password  string
		name      string
		priorTime = 5 * time.Minute
		interval  = 10 * time.Minute
		debug     bool
	)

	godotenv.Load(filepath.Join(dir, EnvFileName))

	account, _ = os.LookupEnv("ACCOUNT")
	if account == "" {
		return nil, errors.New("ACCOUNT is unset")
	}

	password, _ = os.LookupEnv("PASSWORD")
	if password == "" {
		return nil, errors.New("PASSWORD is unset")
	}

	name, _ = os.LookupEnv("NAME")
	if name == "" {
		return nil, errors.New("NAME is unset")
	}

	debugVal, _ := os.LookupEnv("APP_DEBUG")
	debugVal = strings.TrimSpace(debugVal)
	if debugVal != "" && debugVal != "false" && debugVal != "0" {
		debug = true
	} else {
		debug = false
	}

	return &App{
		Account:        account,
		Password:       password,
		Name:           name,
		PriorTime:      priorTime,
		UpdateInterval: interval,
		Debug:          debug,
	}, nil
}
