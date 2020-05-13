package db

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/config"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketName = "cache"
)

type Record struct {
	expiry time.Time
}

type callback func(err error, value *Record)

type Cache struct {
	config config.Configuration
	db     *bolt.DB
}

func (cache *Cache) Get(jti string, cb callback) error {
	handler := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		bytes := bucket.Get([]byte(jti))

		value := &Record{}
		err := json.Unmarshal(bytes, value)

		if err != nil {
			cb(err, nil)
			return err
		}

		cb(nil, value)
		return nil
	}

	err := cache.db.View(handler)
	if err != nil {
		cb(err, nil)
	}
	return nil
}

// Init brings up the embedded BoltDB
func (cache *Cache) Init(config config.Configuration) {
	db, err := bolt.Open("oauth-revokerd.db", 0666, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Info("Database started")
	cache.db = db
}
