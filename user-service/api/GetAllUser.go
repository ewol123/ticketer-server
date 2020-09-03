package api

import (
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type GetAllUserResponseModel struct {
	count int
	rows *[]user.User
}

func (h *handler) GetAllUser(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
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


	getAllRequestModel := user.GetAllUserRequestModel{
		Page:        page,
		RowsPerPage: rowsPerPage,
		SortBy:      sortBy,
		Descending:  descending,
		Filter:      filter,
	}


	res, err := h.userService.GetAllUser(&getAllRequestModel)
	if err != nil {
		if errors.Cause(err) == user.ErrUserInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if errors.Cause(err) == user.ErrRequestInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if errors.Cause(err) == user.ErrUserNotFound {
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
