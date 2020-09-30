package routes

import (
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/fatih/structs"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// GetTicketAdmin : get a ticket for admin
func (h *handler) GetTicketAdmin(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	contentType := r.Header.Get("Content-Type")

	model := ticket.GetTicketRequestModelAdmin{Id: id}


	res, err := h.ticketService.GetTicketAdmin(&model)
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

	newMap := structs.Map(res)


	responseBody, err := h.serializer(contentType).Encode(&newMap)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	setupResponse(w, contentType, responseBody, http.StatusOK)
}