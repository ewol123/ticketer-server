package routes

import (
	"github.com/ewol123/ticketer-server/user-service/serializer/mapdecoder"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userReq, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	updateReqModel := &user.UpdateUserRequestModel{}
	err = mapdecoder.Decode(*userReq, &updateReqModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	updateReqModel.Id = id

	err = h.userService.UpdateUser(updateReqModel)
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

	setupResponse(w, contentType, []byte{}, http.StatusNoContent)
}