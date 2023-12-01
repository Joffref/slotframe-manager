package api

import (
	"fmt"
	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/mux"
	"log/slog"
)

type APIConfig struct {
	Address string `mapstructure:"address"`
}

type API struct {
	cfg      *APIConfig
	router   *mux.Router
	handlers Handler
}

func NewAPI(config *APIConfig, handlers Handler) (*API, error) {
	slog.Debug("Creating API")
	a := &API{
		cfg:      config,
		handlers: handlers,
	}
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	err := r.Handle("/register/{parentId}/{id}/{etx}", mux.HandlerFunc(a.Register))
	if err != nil {
		return nil, err
	}
	a.router = r
	slog.Debug("API created")
	return a, nil
}

func (a *API) Run() error {
	slog.Debug(fmt.Sprintf("Starting API on %s", a.cfg.Address))
	return coap.ListenAndServe("udp", a.cfg.Address, a.router)
}
