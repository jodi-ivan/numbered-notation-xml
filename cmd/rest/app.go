package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	lab "github.com/jodi-ivan/numbered-notation-xml/cmd/rest/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/decorator"
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

	repo := repository.New(context.Background(), db)

	usecaseMod := usecase.New(cfg, repo, renderer.NewRenderer())

	httpRender := adapter.New(
		decorator.WithVariantRedirect(repo)(usecaseMod),
	)

	ws.Register("GET", "/kidung-jemaat/render/:number", httpRender)
	//TODO: make the path root as config
	ws.RegisterStatic("/internal/lab/*filepath", "./files/var/www/html/")
	ws.RegisterStatic("/assets/fonts/*filepath", "./files/var/www/fonts/")

	ws.Register("POST", "/internal/verse-parser", &lab.LyricParser{})

	ws.Register("POST", "/internal/v2/verse-parser", &lab.LyricParserV2{
		Db: db,
	})
	ws.Register("PUT", "/internal/verse/hymn/:hymn/verse/:verse", &lab.VerseManagement{
		VerseRepo: repo,
	})

	ws.Register("PUT", "/internal/v2/verse/hymn/:number", &lab.VerseManagementV2{
		VerseRepo: repo,
		Db:        db,
	})

	err = ws.Serve(cfg.Webserver.Port)
	if err != nil {
		log.Printf("Failed to start the server. Err: %s", err.Error())
		os.Exit(1)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-sigs)

	ws.Stop()
	os.Exit(0)

}
