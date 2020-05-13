package api

import (
	"io"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/willhackett/oauth-revokerd/app/config"
	"github.com/willhackett/oauth-revokerd/app/db"
)

const (
	errInvalidMethod = "Unsupported method"
	errNotFound      = "Record not found"
)

// API produces the methods for the the REST API
type API struct {
	cache  *db.Cache
	config config.Configuration
	logger *log.Entry
}

func (api *API) handleFilterEndpoint(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Get Filter")
}

func (api *API) handleRevokeEndpoint(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		query := req.URL.Query()

		jti := query.Get("jti")
		if jti == "" {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Missing jti from search query")
			return
		}

		api.cache.Get(jti, func(err error, rec *db.Record) {
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, errNotFound)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			io.WriteString(w, "Found")
		})
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, errInvalidMethod)
		return
	}
}

// Init starts the server
func Init(config config.Configuration, cache *db.Cache) {
	logger := log.WithFields(log.Fields{
		"name": "oauth-revokerd",
	})

	api := API{
		cache,
		config,
		logger,
	}

	http.HandleFunc("/revoke", api.handleRevokeEndpoint)
	http.HandleFunc("/filter", api.handleFilterEndpoint)

	portString := strconv.Itoa(api.config.Port)

	log.Println("Server running at on port" + portString + "")
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
