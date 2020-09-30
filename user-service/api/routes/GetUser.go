package routes

import (
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/fatih/structs"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// Get : get a single user by id
func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {


	id := chi.URLParam(r, "id")
	contentType := r.Header.Get("Content-Type")

	findRequestModel := user.GetUserRequestModel{Id: id}


	res, err := h.userService.GetUser(&findRequestModel)
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
