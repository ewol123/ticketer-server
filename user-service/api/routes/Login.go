package routes

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ewol123/ticketer-server/user-service/serializer/mapdecoder"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
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

	loginReqModel := &user.LoginRequestModel{}
	err = mapdecoder.Decode(*userReq, &loginReqModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	loginRes, err := h.userService.Login(loginReqModel)
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

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	claims := jwt.MapClaims{}
	for _, role := range loginRes.Roles {
		claims[role.Name] = true
	}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Hour * 24))
	_, tokenString, err := tokenAuth.Encode(claims)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tokenMap := make(map[string]interface{})
	tokenMap["Token"] = tokenString

	bytes, err := h.serializer(contentType).Encode(&tokenMap)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	setupResponse(w, contentType, bytes, http.StatusOK)
}
