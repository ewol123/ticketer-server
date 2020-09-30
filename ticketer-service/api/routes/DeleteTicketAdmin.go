package routes

import (
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// DeleteTicketAdmin : delete a ticket
func (h *handler) DeleteTicketAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	contentType := r.Header.Get("Content-Type")

	model := ticket.DeleteTicketRequestModelAdmin{Id: id}


	err := h.ticketService.DeleteTicketAdmin(&model)
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


	setupResponse(w, contentType, []byte{}, http.StatusOK)
}