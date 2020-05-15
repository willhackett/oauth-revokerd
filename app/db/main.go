package db

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/config"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketName = "REVOCATIONS"
)

type Record struct {
	ExpiresIn int       `json:"expires_in"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetCallback func(err error, value *Record)

type PutCallback func(err error)

type Cache struct {
	config config.Configuration
	db     *bolt.DB
}

func (cache *Cache) Get(jti string, cb GetCallback) {
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
}

func (cache *Cache) Put(jti string, rec Record, cb PutCallback) {
	handler := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		value, err := json.Marshal(rec)

		if err != nil {
			cb(err)
			return err
		}

		err = bucket.Put([]byte(jti), value)
		if err != nil {
			cb(err)
			return err
		}

		cb(nil)
		return nil
	}

	err := cache.db.Update(handler)
	if err != nil {
		cb(err)
	}
}

func (cache *Cache) Close() {
	defer cache.db.Close()
}

// Init brings up the embedded BoltDB
func Init(config config.Configuration) *Cache {
	var err error
	cache := new(Cache)

	cache.db, err = bolt.Open("oauth-revokerd.db", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("Failed to init database", err)
	}

	err = cache.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			log.Fatal("could not create revocations bucket", err)
			return err
		}
		return nil
	})
	log.Info("Database initialised")
	// defer cache.db.Close()
	return cache
}
