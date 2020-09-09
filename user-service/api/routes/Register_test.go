package routes

import (
	"github.com/ewol123/ticketer-server/user-service/hack"
	"github.com/ewol123/ticketer-server/user-service/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Register)

	req, err := http.NewRequest("POST", "/user/v1/register",
		strings.NewReader(`{"Email":"something.com"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Register test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestRegister(t *testing.T) {

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Register)
	req, err := http.NewRequest("GET", "/user/v1/register",
		strings.NewReader(`{"Email":"info@yourcompanyemail.com","Password":"test", "FullName":"Jack Black"}`))
	if err != nil {
		t.Errorf("Register test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusCreated)
	}

	hack.TearDown()
}
