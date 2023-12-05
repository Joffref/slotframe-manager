package api

import (
	"bytes"
	"encoding/json"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"log"
	"strconv"
)

//go:generate mockgen -package api -destination=./mock_api_handler.go . Handler
type Handler interface {
	Register(parentId string, id string, etx int, input Slots) (*Slots, error)
	Version() int
}

func (a *API) Register(w mux.ResponseWriter, r *mux.Message) {
	vars := r.RouteParams.Vars
	etx, err := strconv.Atoi(vars["etx"])
	if err != nil {
		log.Printf("cannot parse etx: %v", err)
		err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("cannot parse etx")))
		if err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	var registerReponse Slots
	json.NewDecoder(r.Body()).Decode(&registerReponse)
	register, err := a.handlers.Register(vars["parentId"], vars["id"], etx, registerReponse)
	if err != nil {
		log.Printf("cannot register: %v", err)
		switch err.(type) {
		case ErrorParentNodeDoesNotExist:
			err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(err.Error())))
			if err != nil {
				log.Printf("cannot set response: %v", err)
			}
		}
		err = w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("cannot register")))
		if err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	resp, err := json.Marshal(register)
	if err != nil {
		log.Printf("cannot marshal response: %v", err)
		err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("cannot marshal response")))
		if err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	err = w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(resp))
	if err != nil {
		log.Printf("cannot set response: %v", err)
		return
	}
}

func (a *API) Version(w mux.ResponseWriter, r *mux.Message) {
	version := FrameVersion{
		Version: a.handlers.Version(),
	}
	resp, err := json.Marshal(version)
	if err != nil {
		log.Printf("cannot marshal response: %v", err)
		err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("cannot marshal response")))
		if err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(resp))
	if err != nil {
		log.Printf("cannot set response: %v", err)
		return
	}
	return
}
