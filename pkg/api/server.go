package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"io/ioutil"

	"gopkg.in/yaml.v2"

	service "github.com/nathanagood/tapinator/internal/svc"
)

// TapAPIServer is the API server for the tapinator
type TapAPIServer struct {
	router *mux.Router
	svc    service.TapServicer
}

// YamlTapRepository is the repository for the tap list in YAML files.
type YamlTapRepository struct {
	service.TapReader
	service.TapWriter
	filePath string
}

func (r *YamlTapRepository) Write(taps []service.Tap) error {
	log.Debug().Msgf("Saving tap now to file: %", r.filePath)
	file, err := os.Create(r.filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(taps)
	_, err = io.WriteString(file, string(data))
	return err
}

func (r *YamlTapRepository) Read() ([]service.Tap, error) {
	tapList := []service.Tap{}
	yamlFile, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &tapList)
	return tapList, err
}

// NewYamlTapRepository creates a new tap repository for YAML files.
func NewYamlTapRepository(path string) *YamlTapRepository {
	return &YamlTapRepository{
		filePath: path,
	}
}

// NewTapAPIServer creates a new version of the TapAPIServer
func NewTapAPIServer() *TapAPIServer {
	// TODO: will probably get this based on a configuration value that
	// is in a configuration that is passed-in.
	yaml := NewYamlTapRepository("/tmp/taps.yml")

	return &TapAPIServer{
		router: mux.NewRouter(),
		svc: service.NewTapService(
			yaml,
			yaml,
		),
	}
}

// Serve up the API
func (api *TapAPIServer) Serve() {
	log.Debug().Msg("Starting the tap server routes...")
	api.router.HandleFunc("/api/taps", api.getTaps).Methods("GET")
	api.router.HandleFunc("/api/taps", api.saveTap).Methods("POST")

	log.Fatal().Err(http.ListenAndServe(":8080", api.router))
}

func (api *TapAPIServer) getTaps(w http.ResponseWriter, r *http.Request) {
	tapList, err := api.svc.FindAll()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(tapList)
	}
}

func (api *TapAPIServer) saveTap(w http.ResponseWriter, r *http.Request) {
	tap := service.NewTap()
	err := json.NewDecoder(r.Body).Decode(tap)
	if err != nil {
		log.Err(err).Msg("Error while processing body")
		json.NewEncoder(w).Encode(err)
	} else {
		log.Debug().Msgf("Found tap: %s", tap.Name)
		result, err := api.svc.Save(*tap)
		if err != nil {
			log.Err(err).Msg("Error while saving taps")
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(result)
		}
	}
}
