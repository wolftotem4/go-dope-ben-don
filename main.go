package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen2brain/dlgs"
	"github.com/google/logger"
	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/config"
	"github.com/wolftotem4/go-dope-ben-don/internal/app"
	"github.com/wolftotem4/go-dope-ben-don/internal/conducts"
	"github.com/wolftotem4/go-dope-ben-don/internal/gui"
	typesjson "github.com/wolftotem4/go-dope-ben-don/internal/types/json"
)

const (
	appName = "go-dope-ben-don"
)

//go:embed data/* .env.example assets/information.png
var embedFS embed.FS

var performLogout = flag.Bool("logout", false, "Perform logout")
var account = flag.String("account", "", "Provide a account to login dinbendon.net")
var password = flag.String("password", "", "Provide a password to login dinbendon.net")
var name = flag.String("name", "", "Provide a name for matching ordered items")

var Notify chan typesjson.ProgressItem

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {
	ctx := context.Background()

	flag.Parse()
	gob.Register(http.Cookie{})

	dir, err := app.GetConfigDir(appName)
	if err != nil {
		handleFatal(err)
	}

	// generate '.env' if necessary
	if err := generateConfigFile(dir); err != nil {
		handleFatal(err)
	}

	// generate assets
	if err := copyAssets(dir); err != nil {
		handleFatal(err)
	}

	// init error logger
	lf, err := os.OpenFile(filepath.Join(dir, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		handleFatal(errors.Wrap(err, "Failed to open log file"))
	}
	defer lf.Close()

	defer logger.Init("go-dope-ben-don", false, true, lf).Close()

	// init config & client
	app, err := app.Register(dir)
	if err != nil {
		handleFatal(err)
	}

	applyFlagsToConfig(app.Config)

	// init database & load Cookies
	if err := boot(ctx, app); err != nil {
		handleFatal(err)
	}

	if *performLogout {
		log.Println("Logging out...")
		if err := conducts.Logout(ctx, app); err != nil {
			handleFatal(err)
		}
		log.Println("Logged out.")

		if err := app.RealClient.Store.Save(ctx); err != nil {
			handleFatal(err)
		}
	} else {
		g, err := gui.New(dir, appName)
		if err != nil {
			handleFatal(err)
		}
		defer g.Close()

		g.HandleSignals()
		g.Build()

		go func() {
			// must refresh session if it's the first API request
			var firstReq = true

		job:
			for {
				log.Println("Getting prompt items...")
				promptItems, err := conducts.GetPromptItems(ctx, app, firstReq)
				if err != nil {
					handleError(err)
				} else {
					log.Printf("Got prompt items. (%d found)\n", len(*promptItems))

					for _, item := range *promptItems {
						u, err := app.RealClient.Endpoint(item.GetPath())
						if err != nil {
							handleFatal(err)
						}

						if err := g.Notify(gui.PromptItemInfo{
							Name: item.ShopName,
							Url:  u.String(),
						}); err != nil {
							handleFatal(err)
						}
					}
				}
				firstReq = false

				if err := app.RealClient.Store.Save(ctx); err != nil {
					handleFatal(err)
				}

				sleep := time.After(app.Config.PriorTime)

			sleep:
				for {
					select {
					case <-g.Quit:
						break job
					case err := <-g.EventErr:
						handleError(err.Cause())
					case <-sleep:
						break sleep
					}
				}
			}
		}()

		g.Wait()
	}
}

func boot(ctx context.Context, app *app.App) error {
	if err := createTable(ctx, app.DB.DB); err != nil {
		return errors.WithStack(err)
	}

	if err := app.RealClient.Load(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func createTable(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	s, err := embedFS.ReadFile("data/create_tables.sql")
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.ExecContext(ctx, string(s))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func applyFlagsToConfig(config *config.App) {
	if *account != "" {
		config.Account = *account
	}

	if *password != "" {
		config.Password = *password
	}

	if *name != "" {
		config.Name = *name
	}
}

func generateConfigFile(dir string) error {
	if hasAllCLIUserData() {
		return nil
	}

	env := filepath.Join(dir, config.EnvFileName)

	// check if file exists
	if _, err := os.Stat(env); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return errors.WithStack(err)
	}

	// generate '.env' file
	if err := generateFile(".env.example", env); err != nil {
		return errors.WithStack(err)
	}

	// inform user to edit '.env' file
	dlgs.Info("夠多便當", fmt.Sprintf("請編輯 .env 檔案，再重新執行 (%s)", env))

	os.Exit(0)

	return nil
}

func hasAllCLIUserData() bool {
	return *account != "" && *password != "" && *name != ""
}

func handleError(err error) {
	if stackErr, ok := err.(stackTracer); ok {
		var buffer strings.Builder
		for _, f := range stackErr.StackTrace() {
			buffer.WriteString(fmt.Sprintf("%+s:%d\n", f, f))
		}
		logger.Error(buffer.String())
	} else {
		logger.Error(err)
	}
}

func handleFatal(err error) {
	handleError(err)
	dlgs.Error("夠多便當", err.Error())
	os.Exit(1)
}

func copyAssets(dir string) error {
	assetsDir := filepath.Join(dir, "assets")
	if err := os.MkdirAll(assetsDir, os.ModePerm); err != nil {
		return errors.WithStack(err)
	}

	icon := filepath.Join(assetsDir, "information.png")

	// check if file exists
	if _, err := os.Stat(icon); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return errors.WithStack(err)
	}

	if err := generateFile("assets/information.png", icon); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func generateFile(embed string, dst string) error {
	srcFS, err := embedFS.Open(embed)
	if err != nil {
		return errors.WithStack(err)
	}
	defer srcFS.Close()

	dstFS, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer dstFS.Close()

	if _, err := io.Copy(dstFS, srcFS); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
