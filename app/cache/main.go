package database

import (
	"github.com/willhackett/oauth-revokerd/app/config"

	badger "github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
)

type Cache struct {
	config config.Configuration
	db     *badger.DB
}

// Init brings up the etcd database
func (cache *Cache) Init(config config.Configuration) {
	var err error

	cache.db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer cache.store.Close()
}
