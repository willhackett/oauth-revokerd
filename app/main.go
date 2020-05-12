package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/api"
	"github.com/willhackett/oauth-revokerd/app/config"
)

// Init starts the server
func Init() {
	log.SetFormatter(&log.JSONFormatter{})

	config, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", err)
	}

	api.Init(config)
}
