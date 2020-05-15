package api

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/willhackett/oauth-revokerd/app/config"
	"github.com/willhackett/oauth-revokerd/app/db"
)

const (
	errInvalidMethod          = "Unsupported method"
	errNotFound               = "Record not found"
	errUnableToProcessPayload = "Unable to process payload"
)

// API produces the methods for the the REST API
type API struct {
	cache  *db.Cache
	config config.Configuration
	logger *log.Entry
}

func (api *API) handlePostRevocation(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		api.resolve(req, w, http.StatusBadRequest, "Cannot parse form body")
		return
	}

	form := req.Form
	jti := form.Get("jti")
	expiresIn, err := strconv.Atoi(form.Get("expires_in"))

	if err != nil {
		api.resolve(req, w, http.StatusBadRequest, "`expires_in` must be a number of seconds until the record should expire")
		return
	}

	if jti == "" {
		api.resolve(req, w, http.StatusBadRequest, "Body must include `jti` and `expires_in` values")
		return
	}

	body := db.Record{
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
		ExpiresIn: expiresIn,
	}

	api.cache.Put(jti, body, func(err error) {
		if err != nil {
			log.Fatal("Failed to write "+jti+" to database", err)
			api.resolve(req, w, http.StatusInternalServerError, "Failed to write resource")
			return
		}

		api.resolve(req, w, http.StatusCreated, "Created")
	})
	return
}

func (api *API) handleGetRevocation(w http.ResponseWriter, req *http.Request) {
	jti := req.URL.Path[len("/revocations/"):]

	if jti == "" {
		api.resolve(req, w, http.StatusBadRequest, "Missing `jti` from path")
		return
	}

	api.cache.Get(jti, func(err error, rec *db.Record) {
		if err != nil {
			api.resolve(req, w, http.StatusNotFound, errNotFound)
			return
		}

		api.resolve(req, w, http.StatusNoContent, "`jti` exists in store")
	})
}

func (api *API) handlePutRevocation(w http.ResponseWriter, req *http.Request) {

	w.WriteHeader(http.StatusNoContent)
	io.WriteString(w, "Hello")

}

func (api *API) handleDeleteRevocation(w http.ResponseWriter, req *http.Request) {

	w.WriteHeader(http.StatusNoContent)
	io.WriteString(w, "Hello")

}

func (api *API) handleGetFilter(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Get Filter")
}

func (api *API) log(req *http.Request, status int, message string) {
	fields := log.Fields{
		"path":         req.URL.Path,
		"x-request-id": req.Header.Get("x-request-id"),
		"status":       status,
		"message":      message,
	}

	log := api.logger.WithFields(fields)

	switch true {
	case status >= 500:
		log.Error()
		return
	case status >= 400:
		log.Warn()
		return
	default:
		log.Info()
		return
	}
}

func (api *API) resolve(req *http.Request, w http.ResponseWriter, status int, message string) {
	api.log(req, status, message)
	w.WriteHeader(status)
	io.WriteString(w, message)
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

	router := mux.NewRouter()

	router.HandleFunc("/filter", api.handleGetFilter).Methods("GET")
	router.HandleFunc("/revocations", api.handlePostRevocation).Methods("POST")
	router.HandleFunc("/revocations/{id}", api.handleGetRevocation).Methods("GET")
	router.HandleFunc("/revocations/{id}", api.handlePutRevocation).Methods("PUT")
	router.HandleFunc("/revocations/{id}", api.handleDeleteRevocation).Methods("DELETE")

	portString := strconv.Itoa(api.config.Port)

	http.Handle("/", router)

	log.Println("API available at http://127.0.0.1:" + portString + "/")
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
