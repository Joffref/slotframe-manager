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
	r := newRouter()
	if err := a.registerHandlers(r); err != nil {
		return nil, err
	}
	a.router = r
	slog.Debug("API created")
	return a, nil
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	return r
}

func (a *API) registerHandlers(r *mux.Router) error {
	err := r.Handle("/register/{parentId}/{id}/{etx}", mux.HandlerFunc(a.Register))
	if err != nil {
		return ErrorUnableToRegisterHandler{
			Pattern: "/register/{parentId}/{id}/{etx}",
			err:     err,
		}
	}
	err = r.Handle("/version", mux.HandlerFunc(a.Version))
	if err != nil {
		return ErrorUnableToRegisterHandler{
			Pattern: "/version",
			err:     err,
		}
	}
	return nil
}

func (a *API) Run() error {
	slog.Debug(fmt.Sprintf("Starting API on %s", a.cfg.Address))
	return coap.ListenAndServe("udp", a.cfg.Address, a.router)
}
