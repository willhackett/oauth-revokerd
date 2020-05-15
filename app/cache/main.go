package cache

import (
	"context"
	"errors"
	"time"

	"github.com/buraksezer/olric"
	olricCfg "github.com/buraksezer/olric/config"
	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/config"
)

var (
	bucketName         = "revocations"
	errNotFound        = errors.New("Record not found")
	errUnpackingRecord = errors.New("Error unpacking record")
	errPackingRecord   = errors.New("Error packing record")
)

type Cache struct {
	db *olric.Olric
	dm *olric.DMap
}

// Put creates a new record or overrides it
func (cache *Cache) Put(jti string, expiresIn time.Duration) error {

	expiresAt := time.Now().Add(expiresIn).Unix()

	err := cache.dm.PutEx(jti, expiresAt, expiresIn)
	if err != nil {
		return errPackingRecord
	}

	return nil
}

// Get retrieves a record and its expiry
func (cache *Cache) Get(jti string) (time.Time, error) {
	data, err := cache.dm.Get(jti)
	if err != nil {
		return time.Now(), errNotFound
	}

	timeInt, ok := data.(int64)
	if !ok {
		return time.Now(), errUnpackingRecord
	}

	expiresAt := time.Unix(timeInt, 0)

	return expiresAt, nil
}

// Init brings up the embedded store
func Init(config config.Configuration) *Cache {
	defaults := olricCfg.New("local")

	defaults.BindPort = config.CachePort

	ctx, cancel := context.WithCancel(context.Background())

	defaults.Started = func() {
		defer cancel()
		log.Println("[INFO] Olric is ready to accept connections")
	}

	db, err := olric.New(defaults)

	if err != nil {
		log.Fatalf("Failed to create Olric instance: %v", err)
	}

	go func() {
		// Call Start at background. It's a blocker call.
		err = db.Start()
		if err != nil {
			log.Fatalf("olric.Start returned an error: %v", err)
		}
	}()

	<-ctx.Done()

	dm, err := db.NewDMap(bucketName)
	if err != nil {
		log.Fatalf("olric.NewDMap returned an error: %v", err)
	}

	return &Cache{
		db,
		dm,
	}
}
