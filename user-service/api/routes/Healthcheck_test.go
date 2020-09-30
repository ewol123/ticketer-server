package routes

import (
	"github.com/ewol123/ticketer-server/user-service/user"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestHealthcheck(t *testing.T) {
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	req, err := http.NewRequest("GET", "/user/v1/healthcheck", nil)
	if err != nil {
		t.Errorf("Healtcheck test failed, error: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Healthcheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusOK)
	}


}