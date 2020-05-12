package api

import (
	"io"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/willhackett/oauth-revokerd/app/config"
)

type API struct {
	config config.Configuration
	logger *log.Entry
}

func (api *API) handlePostRevoke(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Post Revoke")
}

func (api *API) handleGetFilter(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Get Filter")
}

func (api *API) handleGetToken(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Get Filter")
}

// Init starts the server
func Init(config config.Configuration) {
	logger := log.WithFields(log.Fields{
		"name": "oauth-revokerd",
	})

	api := API{
		config,
		logger,
	}

	http.HandleFunc("/v1/revoke", api.handlePostRevoke)
	http.HandleFunc("/v1/filter", api.handleGetFilter)
	http.HandleFunc("/v1/token", api.handleGetToken)

	portString := strconv.Itoa(api.config.Port)
	log.Println("Server running at http://localhost:" + portString + "")
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
