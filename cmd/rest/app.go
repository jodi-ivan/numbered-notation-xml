package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/renderer"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/svc/usecase"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
	"github.com/jodi-ivan/numbered-notation-xml/utils/storage"
	"github.com/jodi-ivan/numbered-notation-xml/utils/webserver"
)

func main() {
	env := "development"
	cfg, err := config.InitConfig(env)
	if err != nil {
		log.Fatalf("failed to load config, err : %s", err.Error())
	}
	ws, err := webserver.InitWebserver()
	if err != nil {
		log.Fatalf("[Webserver] Failed to initilize the webserver. err : %s", err.Error())
		return
	}

	db, err := storage.NewStorage(context.Background(), cfg.SQLite.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to storage: %s", err.Error())
		return
	}

	repo := repository.New(context.Background(), db, cfg)

	usecaseMod := usecase.New(cfg, repo, renderer.NewRenderer())

	httpRender := adapter.New(usecaseMod)

	ws.Register("GET", "/kidung-jemaat/render/:number", httpRender)

	ws.Serve(cfg.Webserver.Port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-sigs)

	ws.Stop()

}
