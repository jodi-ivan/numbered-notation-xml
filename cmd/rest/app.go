package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jodi-ivan/numbered-notation-xml/adapter"
	"github.com/jodi-ivan/numbered-notation-xml/utils/config"
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
	}

	ws.Register("GET", "/test", &adapter.TrySVG{})

	ws.Serve(cfg.Webserver.Port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-sigs)

	ws.Stop()

}
