package app

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/config"
	"github.com/wolftotem4/go-dope-ben-don/internal/client"
)

type App struct {
	Config     *config.App
	DB         *sqlx.DB
	RealClient *client.HttpClient
	Client     client.Client
}

func Register() (*App, error) {
	var (
		cfig    *config.App
		db      *sqlx.DB
		real    *client.HttpClient
		client_ client.Client
		err     error
	)

	cfig, err = config.LoadConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db, err = sqlx.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	real, err = client.NewClient(client.NewDatabaseStore(db))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if cfig.Debug {
		client_ = client.WrapLogger(real)
	} else {
		client_ = real
	}

	return &App{
		Config:     cfig,
		DB:         db,
		RealClient: real,
		Client:     client_,
	}, nil
}
