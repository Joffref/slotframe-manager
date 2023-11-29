package api

import (
	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/mux"
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
	return a, nil
}

func (a *API) Run() error {
	return coap.ListenAndServe("udp", a.cfg.Address, a.router)
}
