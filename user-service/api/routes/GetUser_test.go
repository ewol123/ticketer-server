package routes

import (
	"context"
	"encoding/json"
	"github.com/ewol123/ticketer-server/user-service/hack"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetUser)
	req, err := http.NewRequest("GET", "/user/v1/asd", nil)
	if err != nil {
		t.Errorf("GetUser test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetUserNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetUser)
	req, err := http.NewRequest("GET", "/user/v1/", nil)
	if err != nil {
		t.Errorf("GetUser test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "c655b6b9-3956-4ee9-910a-2560e8e49d6e")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestGetUserFound(t *testing.T) {

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetUser)
	req, err := http.NewRequest("GET", "/user/v1/", nil)
	if err != nil {
		t.Errorf("GetUser test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "e66c0b06-ec6c-45e1-8619-27e14c3ed92d")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
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
		t.Errorf("cannot decode response body of GetUser %v", err)
	}

	isEqual := m["Email"] == "peti@test.com" && m["FullName"] == "peti" && m["RegistrationCode"] == ""

	if isEqual != true {
		t.Errorf("handler returned wrong response body")
	} else {
		t.Logf("handler returned correct response body")
	}
	hack.TearDown()

}
