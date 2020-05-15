package app

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/api"
	"github.com/willhackett/oauth-revokerd/app/config"
	"github.com/willhackett/oauth-revokerd/app/db"
)

func cleanup(cache *db.Cache) {
	defer cache.Close()
}

// Init starts the server
func Init() {
	log.SetFormatter(&log.JSONFormatter{})

	config, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", err)
	}

	cache := db.Init(config)

	defer api.Init(config, cache)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(cache)
		os.Exit(1)
	}()
}
