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

type RegisterReponse struct {
	Slots []int `json:"slots"`
}

type Handler interface {
	Register(parentId string, id string, etx int, input RegisterReponse) (*RegisterReponse, error)
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
	var registerReponse RegisterReponse
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
