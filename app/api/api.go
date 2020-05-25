package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/willhackett/oauth-revokerd/app/filter"

	"github.com/willhackett/oauth-revokerd/app/cache"
	"github.com/willhackett/oauth-revokerd/app/config"
)

const (
	errInvalidMethod          = "Unsupported method"
	errNotFound               = "Record not found"
	errUnableToProcessPayload = "Unable to process payload"
)

// API produces the methods for the the REST API
type API struct {
	cache  *cache.Cache
	config config.Configuration
	logger *log.Entry
	filter *filter.Filter
}

func (api *API) handlePostRevocation(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		api.resolve(req, w, http.StatusBadRequest, "Cannot parse form body")
		return
	}

	form := req.Form
	jti := form.Get("jti")
	expiresIn, err := strconv.ParseInt(form.Get("expires_in"), 10, 64)

	if err != nil {
		api.resolve(req, w, http.StatusBadRequest, "`expires_in` must be a number of seconds until the record should expire")
		return
	}

	if jti == "" {
		api.resolve(req, w, http.StatusBadRequest, "Body must include `jti` and `expires_in` values")
		return
	}

	expiresInDuration := time.Duration(expiresIn) * time.Second

	err = api.cache.Put(jti, expiresInDuration)

	if err != nil {
		api.resolve(req, w, http.StatusInternalServerError, "Failed to write resource")
		return
	}

	api.resolve(req, w, http.StatusCreated, "Created")
	return
}

func (api *API) handleGetRevocation(w http.ResponseWriter, req *http.Request) {
	jti := req.URL.Path[len("/revocations/"):]

	if jti == "" {
		api.resolve(req, w, http.StatusBadRequest, "Missing `jti` from path")
		return
	}

	expiresAt, err := api.cache.Get(jti)

	if err != nil {
		api.resolve(req, w, http.StatusNotFound, errNotFound)
		return
	}

	isInFilter := api.filter.Test(jti)

	resStr := "`jti` revocation expires at " + expiresAt.Local().String()

	if isInFilter {
		resStr += "- present in filter."
	}

	api.resolve(req, w, http.StatusNoContent, resStr)
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
	total, err := api.cache.Count()
	if err != nil {
		api.resolve(req, w, http.StatusInternalServerError, fmt.Sprintf("Failed to get count: %s", err.Error()))
		return
	}

	filter := filter.New(uint(total))

	err = api.cache.Query(func(jti string) {
		filter.Add(jti)
	})

	api.filter = filter

	if err != nil {
		api.resolve(req, w, http.StatusInternalServerError, fmt.Sprintf("Failed to populate filter: %s", err.Error()))
		return
	}

	filterJSON, err := filter.MarshalJSON()
	if err != nil {
		api.resolve(req, w, http.StatusInternalServerError, fmt.Sprintf("Failed to export filter: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(filterJSON)
	api.log(req, http.StatusOK, "Filter Supplied")
	if err != nil {
		api.resolve(req, w, http.StatusInternalServerError, fmt.Sprintf("Failed to write response: %s", err.Error()))
		return
	}
}

func (api *API) handleGetCount(w http.ResponseWriter, req *http.Request) {
	count, err := api.cache.Count()
	if err != nil {
		api.resolve(req, w, http.StatusInternalServerError, err.Error())
		return
	}

	api.resolve(req, w, http.StatusOK, strconv.Itoa(count))
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
func Init(config config.Configuration, cache *cache.Cache) {
	logger := log.WithFields(log.Fields{
		"name": "oauth-revokerd",
	})

	filter := filter.New(1)

	api := API{
		cache,
		config,
		logger,
		filter,
	}

	router := mux.NewRouter()

	router.HandleFunc("/filter", api.handleGetFilter).Methods("GET")
	router.HandleFunc("/count", api.handleGetCount).Methods("GET")
	router.HandleFunc("/revocations", api.handlePostRevocation).Methods("POST")
	router.HandleFunc("/revocations/{id}", api.handleGetRevocation).Methods("GET")
	router.HandleFunc("/revocations/{id}", api.handlePutRevocation).Methods("PUT")
	router.HandleFunc("/revocations/{id}", api.handleDeleteRevocation).Methods("DELETE")

	portString := strconv.Itoa(api.config.Port)

	http.Handle("/", router)

	log.Println("API available at http://127.0.0.1:" + portString + "/")
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
