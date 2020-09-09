package routes

import (
	"encoding/json"
	"github.com/ewol123/ticketer-server/user-service/repository/postgres/seed"
	"github.com/ewol123/ticketer-server/user-service/user"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestGetAllUserBadRequest(t *testing.T) {
	seed.Init("../../repository/postgres/seed/seed_test.sql")

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAllUser)
	req, err := http.NewRequest("GET", "/user/v1", nil)
	if err != nil {
		t.Errorf("GetAllUser test failed, error: %v", err)
	}

	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("rowsPerPage", "10")
	q.Add("descending", "false")
	req.URL.RawQuery = q.Encode()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetAllUserFound(t *testing.T){

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAllUser)
	req, err := http.NewRequest("GET", "/user/v1", nil)
	if err != nil {
		t.Errorf("GetAllUser test failed, error: %v", err)
	}

	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("rowsPerPage", "10")
	q.Add("sortBy", "user_id")
	q.Add("descending", "false")
	q.Add("filter", "")
	req.URL.RawQuery = q.Encode()
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
		t.Errorf("cannot decode response body of GetAllUser %v", err)
	}

	if m["Count"] != 0 {
		t.Logf("respone body looks OK %v",m)
	} else {
		t.Errorf("response body incorrect")
	}

	seed.TearDown()
}