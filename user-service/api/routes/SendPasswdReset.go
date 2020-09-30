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

func (h *handler) SendPasswdReset(w http.ResponseWriter, r *http.Request) {
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

	sendPwResetReqModel := &user.SendPasswdResetRequestModel{}
	err = mapdecoder.Decode(*userReq, &sendPwResetReqModel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	storedUser, err := h.userService.SendPasswdReset(sendPwResetReqModel)
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


	from := mail.NewEmail("Your password reset", "gypeti23@gmail.com")
	subject := "Reset your password with the provided code"
	to := mail.NewEmail(storedUser.FullName, storedUser.Email)
	plainTextContent := fmt.Sprintf(`
		please enter the provided code in the application to complete your password reset: %v`,storedUser.ResetPasswordCode)
	htmlContent := fmt.Sprintf(`<strong>
		please enter the provided code in the application to complete your password reset: %v</strong>`,storedUser.ResetPasswordCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	setupResponse(w, contentType, []byte{}, http.StatusOK)
}
