package config

import (
	"errors"
	"os"
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
}

func LoadConfig() (*App, error) {
	godotenv.Load()

	account, _ := os.LookupEnv("ACCOUNT")
	if account == "" {
		return nil, errors.New("ACCOUNT is unset")
	}

	password, _ := os.LookupEnv("PASSWORD")
	if password == "" {
		return nil, errors.New("PASSWORD is unset")
	}

	name, _ := os.LookupEnv("NAME")
	if name == "" {
		return nil, errors.New("NAME is unset")
	}

	return &App{
		Account:   account,
		Password:  password,
		Name:      name,
		PriorTime: 5 * time.Minute,
	}, nil
}
