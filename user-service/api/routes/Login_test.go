package routes

import (
	"encoding/json"
	"github.com/ewol123/ticketer-server/user-service/hack"
	"github.com/ewol123/ticketer-server/user-service/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Login)

	req, err := http.NewRequest("POST", "/user/v1/login",
		strings.NewReader(`{"Email":"something.com"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Login test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestLoginNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Login)
	req, err := http.NewRequest("POST", "/user/v1/login",
		strings.NewReader(`{"Email":"test.fail@test.com","Password":"asd"}`))
	if err != nil {
		t.Errorf("Login test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestLogin(t *testing.T) {

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Login)
	req, err := http.NewRequest("GET", "/user/v1/login",
		strings.NewReader(`{"Email":"test7@test.com","Password":"test"}`))
	if err != nil {
		t.Errorf("Login test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusOK)
	}

	// check response body
	m := make(map[string]interface{})
	if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
		t.Errorf("cannot decode response body of Login %v", err)
	}

	if m["Token"] == "" {
		t.Errorf("handler returned wrong response body")
	} else {
		t.Logf("handler returned correct response body: %v", m)
	}


	hack.TearDown()
}
