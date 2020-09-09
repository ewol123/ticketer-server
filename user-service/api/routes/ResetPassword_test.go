package routes

import (
	"github.com/ewol123/ticketer-server/user-service/hack"
	"github.com/ewol123/ticketer-server/user-service/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResetPasswordBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ResetPassword)

	req, err := http.NewRequest("POST", "/user/v1/reset-password",
		strings.NewReader(`{"Email":"something.com"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("ResetPassword test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestResetPasswordNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ResetPassword)
	req, err := http.NewRequest("POST", "/user/v1/reset-password",
		strings.NewReader(`{"Email":"test.fail@test.com","Password":"asd", "ResetPasswordCode": "23442"}`))
	if err != nil {
		t.Errorf("ResetPassword test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestResetPassword(t *testing.T) {

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ResetPassword)
	req, err := http.NewRequest("GET", "/user/v1/reset-password",
		strings.NewReader(`{"Email":"test8@test.com","Password":"test", "ResetPasswordCode": "123456"}`))
	if err != nil {
		t.Errorf("ResetPassword test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusOK)
	}


	hack.TearDown()
}
