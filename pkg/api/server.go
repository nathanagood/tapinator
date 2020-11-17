package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	service "github.com/nathanagood/tapinator/internal/svc"
)

// TapAPIServer is the API server for the tapinator
type TapAPIServer struct {
	router *mux.Router
	svc    service.TapServicer
}

// NewTapAPIServer creates a new version of the TapAPIServer
func NewTapAPIServer() *TapAPIServer {
	return &TapAPIServer{
		router: mux.NewRouter(),
		svc:    service.NewTapService(),
	}
}

// Serve up the API
func (api *TapAPIServer) Serve() {
	log.Debug().Msg("Starting the tap server routes...")
	api.router.HandleFunc("/api/", api.getTaps).Methods("GET")

	log.Fatal().Err(http.ListenAndServe(":8080", api.router))
}

func (api *TapAPIServer) getTaps(w http.ResponseWriter, r *http.Request) {

	tapList, err := api.svc.GetTaps()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(tapList)
	}
}
