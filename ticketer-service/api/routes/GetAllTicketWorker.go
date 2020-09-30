package routes

import (
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// GetAllTicketWorker : get all ticket for workers
func (h *handler) GetAllTicketWorker(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	requesterId := r.Header.Get("Requester-Id")

	queryParams := r.URL.Query()


	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rowsPerPage, err := strconv.Atoi(queryParams.Get("rowsPerPage"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	descending,err := strconv.ParseBool(queryParams.Get("descending"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sortBy := queryParams.Get("sortBy")
	filter := queryParams.Get("filter")


	model := ticket.GetAllTicketRequestModelWorker{
		Page:        page,
		RowsPerPage: rowsPerPage,
		SortBy:      sortBy,
		Descending:  descending,
		Filter:      filter,
		RequesterId: requesterId,
	}


	res, err := h.ticketService.GetAllTicketWorker(&model)
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