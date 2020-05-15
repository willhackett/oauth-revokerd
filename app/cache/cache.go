package cache

import (
	"context"
	"errors"
	"time"

	"github.com/buraksezer/olric"
	olricCfg "github.com/buraksezer/olric/config"
	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/config"
	"github.com/willhackett/oauth-revokerd/app/discovery"
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
	cfg := olricCfg.New(config.MemberlistConfig)

	disco := &discovery.CloudDiscovery{}

	cfg.LogVerbosity = 6

	if config.DiscoveryProvider != "" {
		cfg.ServiceDiscovery = map[string]interface{}{
			"plugin":   disco,
			"provider": config.DiscoveryProvider,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg.Started = func() {
		defer cancel()
		log.Info("Store is ready to accept connections")
	}

	db, err := olric.New(cfg)

	if err != nil {
		log.Fatal("Failed to create cache instance", err)
	}

	go func() {
		// Call Start at background. It's a blocker call.
		err = db.Start()
		if err != nil {
			log.Fatal("Failed to start cache", err)
		}
	}()

	<-ctx.Done()

	dm, err := db.NewDMap(bucketName)
	if err != nil {
		log.Fatal("Failed to create bucket", err)
	}

	return &Cache{
		db,
		dm,
	}
}
