package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"satellity/internal/configs"
	"satellity/internal/controllers"
	"satellity/internal/durable"
	"satellity/internal/middlewares"
	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/handlers"
	flags "github.com/jessevdk/go-flags"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func startHTTP(db *sql.DB, logger *zap.Logger, port string) error {
	database := durable.WrapDatabase(db)
	router := httptreemux.New()
	controllers.RegisterHanders(router)
	controllers.RegisterRoutes(database, router)

	handler := middlewares.Authenticate(database, router)
	handler = middlewares.Constraint(handler)
	handler = middlewares.Context(handler, render.New())
	handler = middlewares.State(handler)
	handler = middlewares.Logger(handler, durable.NewLogger(logger))
	handler = handlers.ProxyHeaders(handler)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

func main() {
	var options struct {
		Config      string `short:"c" long:"config" description:"Where's the config file place, default ./internal/configs/config.yaml"`
		Environment string `short:"e" long:"environment" default:"development"`
	}
	p := flags.NewParser(&options, flags.Default)
	if _, err := p.Parse(); err != nil {
		log.Panicln(err)
	}

	if options.Config == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Panicln(err)
		}
		back := ".."
		if strings.Contains(dir, "cmd") {
			back = "../.."
		}
		options.Config = path.Join(dir, back, "internal/configs/config.yaml")
	}

	if err := configs.Init(options.Config, options.Environment); err != nil {
		log.Panicln(err)
	}

	config := configs.AppConfig
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     config.Database.User,
		Password: config.Database.Password,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Name:     config.Database.Name,
	})
	defer db.Close()

	logger, err := zap.NewDevelopment()
	if config.Environment == "production" {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Panicln(err)
	}

	if err := startHTTP(db, logger, config.HTTP.Port); err != nil {
		log.Panicln(err)
	}
}
