package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/ewol123/ticketer-server/ticketer-service/serializer/mapdecoder"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/pkg/errors"
)

// CreateTicketUser : create a new ticket
func (h *handler) CreateTicketUser(w http.ResponseWriter, r *http.Request) {
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

	model := &ticket.CreateTicketRequestModelUser{}
	err = mapdecoder.Decode(*request, &model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	model["RequesterId"] = requesterId

	err = h.ticketService.CreateTicketUser(model)
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

	setupResponse(w, contentType, []byte{}, http.StatusCreated)

}
