package api

import (
	js "github.com/ewol123/ticketer-server/user-service/serializer/json"
	"github.com/ewol123/ticketer-server/user-service/user"
	"log"
	"net/http"
)

// UserHandler : UserHandler interface with Get and Post methods
type UserHandler interface {
	GetUser(http.ResponseWriter, *http.Request)
	GetAllUser(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	ConfirmRegistration(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	SendPasswdReset(http.ResponseWriter, *http.Request)
	ResetPassword(http.ResponseWriter, *http.Request)
}


type handler struct {
	userService user.Service
}


func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}


func (h *handler) serializer(contentType string) user.Serializer {
	/*if contentType = "application/x-msgpack" {
		return &ms.User{}
	}*/
	return &js.User{}
}


// NewHandler : returns a new UserHandler
func NewHandler(userService user.Service) UserHandler {
	return &handler{userService: userService}
}