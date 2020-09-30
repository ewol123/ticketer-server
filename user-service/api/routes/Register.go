package routes

import (
	"fmt"
	"github.com/ewol123/ticketer-server/user-service/serializer/mapdecoder"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"io/ioutil"
	"net/http"
	"os"
)

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {

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

	regReqModel := &user.RegisterRequestModel{}
	err = mapdecoder.Decode(*userReq, &regReqModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	storedUser, err := h.userService.Register(regReqModel)
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


	deleteRequestModel := user.DeleteUserRequestModel{Id: storedUser.Id}


	from := mail.NewEmail("Complete your registration", "gypeti23@gmail.com")
	subject := "Activate your registration with the provided code"
	to := mail.NewEmail(regReqModel.FullName, regReqModel.Email)
	plainTextContent := fmt.Sprintf(`
		please enter the provided code in the application to complete your registration: %v`,storedUser.RegistrationCode)
	htmlContent := fmt.Sprintf(`<strong>
		please enter the provided code in the application to complete your registration: %v</strong>`,storedUser.RegistrationCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		err = h.userService.DeleteUser(&deleteRequestModel)
		if err != nil {
			if errors.Cause(err) == user.ErrUserNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	setupResponse(w, contentType, []byte{}, http.StatusCreated)

}