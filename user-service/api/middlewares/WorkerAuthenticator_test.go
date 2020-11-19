package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func TestWorkerAuthenticator(t *testing.T) {
	r := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	r.Use(jwtauth.Verifier(tokenAuth), WorkerAuthenticator)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// sending unauthorized requests
	if status, resp := testRequest(t, ts, "GET", "/", nil, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator without token passed want %v got %v", 401, status)
	}

	h := http.Header{}
	h.Set("Authorization", "BEARER "+newJwtToken([]byte("wrong"), map[string]interface{}{}))
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator with wrong secret passed want %v got %v", 401, status)
	}

	h.Set("Authorization", "BEARER asdf")
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator with wrong token passed want %v got %v", 401, status)
	}

	// wrong token secret and wrong alg
	h.Set("Authorization", "BEARER "+newJwt512Token([]byte("wrong"), map[string]interface{}{}))
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator with wrong token secret and wrong alg passed want %v got %v", 401, status)
	}

	// correct token secret but wrong alg
	h.Set("Authorization", "BEARER "+newJwt512Token([]byte(os.Getenv("JWT_SECRET")), map[string]interface{}{}))
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator with correct token secret but wrong alg passed want %v got %v", 401, status)
	}

	//wrong claim
	adminClaims := jwt.MapClaims{}
	adminClaims["admin"] = true
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(adminClaims), nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator with wrong claim passed want %v got %v", 401, status)
	}

	// sending authorized requests
	userClaims := jwt.MapClaims{}
	userClaims["worker"] = true
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(userClaims), nil); status != 200 || resp != "welcome" {
		t.Errorf(resp)
	} else {
		t.Logf("WorkerAuthenticator with correct claim passed want %v got %v", 200, status)
	}
}
