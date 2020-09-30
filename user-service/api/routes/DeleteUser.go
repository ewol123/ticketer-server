package routes

import (
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	contentType := r.Header.Get("Content-Type")

	deleteRequestModel := user.DeleteUserRequestModel{Id: id}


	err := h.userService.DeleteUser(&deleteRequestModel)
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


	setupResponse(w, contentType, []byte{}, http.StatusOK)
}
