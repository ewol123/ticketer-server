package routes

import (
	"github.com/ewol123/ticketer-server/ticketer-service/serializer/mapdecoder"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// UpdateTicketAdmin : update a ticket for admin users
func (h *handler) UpdateTicketAdmin(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	contentType := r.Header.Get("Content-Type")
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

	model := &ticket.UpdateTicketRequestModelAdmin{}
	err = mapdecoder.Decode(*request, &model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	model.Id = id

	err = h.ticketService.UpdateTicketAdmin(model)
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

	setupResponse(w, contentType, []byte{}, http.StatusNoContent)
}