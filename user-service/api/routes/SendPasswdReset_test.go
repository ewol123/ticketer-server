package routes

import (
	"github.com/ewol123/ticketer-server/user-service/hack"
	"github.com/ewol123/ticketer-server/user-service/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendPasswdResetBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SendPasswdReset)

	req, err := http.NewRequest("POST", "/user/v1/send-passwd-reset",
		strings.NewReader(`{"Email":"something.com"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("SendPasswdReset test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSendPasswdResetNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SendPasswdReset)
	req, err := http.NewRequest("POST", "/user/v1/send-passwd-reset",
		strings.NewReader(`{"Email":"test.fail@test.com"}`))
	if err != nil {
		t.Errorf("SendPasswdReset test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestSendPasswdReset(t *testing.T) {

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SendPasswdReset)
	req, err := http.NewRequest("GET", "/user/v1/send-passwd-reset",
		strings.NewReader(`{"Email":"test9@test.com"}`))
	if err != nil {
		t.Errorf("SendPasswdReset test failed, error: %v", err)
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
