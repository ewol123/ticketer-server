package routes

import (
	"github.com/ewol123/ticketer-server/ticketer-service/serializer/mapdecoder"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// SyncTicketWorker : ticket sync between client and server
func (h *handler) SyncTicketWorker(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requesterId := r.Header.Get("Requester-Id")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	request, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	model := &ticket.SyncTicketRequestModelWorker{}
	err = mapdecoder.Decode(*request,&model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	model.RequesterId = requesterId

	serviceResponse, err := h.ticketService.SyncTicketWorker(model)
	if err != nil {
		if errors.Cause(err) == ticket.ErrTicketInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if errors.Cause(err) == ticket.ErrRequestInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if errors.Cause(err) == ticket.ErrTicketNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newMap := structs.Map(serviceResponse)

	responseBody, err := h.serializer(contentType).Encode(&newMap)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	setupResponse(w, contentType, responseBody, http.StatusOK)
}