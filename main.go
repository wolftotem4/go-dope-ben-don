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
	"time"

	"github.com/gen2brain/beeep"
	"github.com/pkg/errors"
	"github.com/wolftotem4/go-dope-ben-don/internal/app"
	"github.com/wolftotem4/go-dope-ben-don/internal/conducts"
	typesjson "github.com/wolftotem4/go-dope-ben-don/internal/types/json"
)

//go:embed data/* .env.example
var embedFS embed.FS

var performLogout = flag.Bool("logout", false, "Perform logout")
var generateEnv = flag.Bool("generate", false, "Generate .env file")

func main() {
	ctx := context.Background()

	flag.Parse()
	gob.Register(http.Cookie{})

	if *generateEnv {
		log.Println("Generating '.env'...")

		if _, err := os.Stat(".env"); err == nil {
			log.Println("'.env' already exists.")
		} else if errors.Is(err, os.ErrNotExist) {
			src, err := embedFS.Open(".env.example")
			if err != nil {
				log.Fatal(err)
			}
			defer src.Close()

			dst, err := os.Create(".env")
			if err != nil {
				log.Fatal(err)
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				log.Fatal(err)
			}

			log.Println("'.env' is generated.")
		} else {
			log.Fatal(err)
		}
	} else {
		app, err := app.Register()
		if err != nil {
			log.Fatal(err)
		}

		if err := boot(ctx, app); err != nil {
			log.Fatal(err)
		}

		if *performLogout {
			log.Println("Logging out...")
			if err := conducts.Logout(ctx, app); err != nil {
				log.Fatal(err)
			}
			log.Println("Logged out.")

			if err := app.RealClient.Store.Save(ctx); err != nil {
				log.Fatal(err)
			}
		} else {
			var first = true

			for {
				log.Println("Getting prompt items...")
				promptItems, err := conducts.GetPromptItems(ctx, app, first)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Got prompt items. (%d found)\n", len(*promptItems))

					if len(*promptItems) > 0 {
						if err := notify(promptItems); err != nil {
							log.Fatal(err)
						}
					}
				}
				first = false

				if err := app.RealClient.Store.Save(ctx); err != nil {
					log.Fatal(err)
				}

				time.Sleep(app.Config.PriorTime)
			}
		}
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

func notify(items *[]typesjson.ProgressItem) error {
	message := fmt.Sprintf("%d 個項目快過期...", len(*items))

	if err := beeep.Notify("夠多便當", message, "assets/information.png"); err != nil {
		return err
	}
	return nil
}
