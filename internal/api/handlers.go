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
	var registerResponse Slots
	if err = json.NewDecoder(r.Body()).Decode(&registerResponse); err != nil {
		if err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(err.Error()))); err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	register, err := a.handlers.Register(vars["parentId"], vars["id"], etx, registerResponse)
	if err != nil {
		log.Printf("cannot register: %v", err)
		switch err.(type) {
		case ErrorParentNodeDoesNotExist:
			if err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(err.Error()))); err != nil {
				log.Printf("cannot set response: %v", err)
			}
		default:
			if err = w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte("cannot register"))); err != nil {
				log.Printf("cannot set response: %v", err)
			}
		}
		return
	}
	resp, err := json.Marshal(register)
	if err != nil {
		log.Printf("cannot marshal response: %v", err)
		if err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("cannot marshal response"))); err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	if err = w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(resp)); err != nil {
		log.Printf("cannot set response: %v", err)
		return
	}
}

func (a *API) Version(w mux.ResponseWriter, r *mux.Message) {
	version := FrameVersion{
		Version: a.handlers.Version(),
	}
	resp, err := json.Marshal(&version)
	if err != nil {
		log.Printf("cannot marshal response: %v", err)
		if err = w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("cannot marshal response"))); err != nil {
			log.Printf("cannot set response: %v", err)
		}
		return
	}
	if err = w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(resp)); err != nil {
		log.Printf("cannot set response: %v", err)
		return
	}
	return
}
