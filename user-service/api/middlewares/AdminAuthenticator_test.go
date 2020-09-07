package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, header http.Header, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func newJwtToken(secret []byte, claims ...jwt.MapClaims) string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	if len(claims) > 0 {
		token.Claims = claims[0]
	}
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
	}
	return tokenStr
}

func newJwt512Token(secret []byte, claims ...jwt.MapClaims) string {
	// use-case: when token is signed with a different alg than expected
	token := jwt.New(jwt.GetSigningMethod("HS512"))
	if len(claims) > 0 {
		token.Claims = claims[0]
	}
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
	}
	return tokenStr
}

func newAuthHeader(claims ...jwt.MapClaims) http.Header {
	h := http.Header{}
	h.Set("Authorization", "BEARER "+newJwtToken([]byte(os.Getenv("JWT_SECRET")), claims...))
	return h
}

func TestAdminAuthenticator(t *testing.T){
	r := chi.NewRouter()
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	r.Use(jwtauth.Verifier(tokenAuth), AdminAuthenticator)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// sending unauthorized requests
	if status, resp := testRequest(t, ts, "GET", "/", nil, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator without token passed want %v got %v", 401,status)
	}

	h := http.Header{}
	h.Set("Authorization", "BEARER "+newJwtToken([]byte("wrong"), map[string]interface{}{}))
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator with wrong secret passed want %v got %v", 401, status)
	}

	h.Set("Authorization", "BEARER asdf")
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator with wrong token passed want %v got %v",401, status)
	}

	// wrong token secret and wrong alg
	h.Set("Authorization", "BEARER "+newJwt512Token([]byte("wrong"), map[string]interface{}{}))
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator with wrong token secret and wrong alg passed want %v got %v",401,status)
	}

	// correct token secret but wrong alg
	h.Set("Authorization", "BEARER "+newJwt512Token([]byte(os.Getenv("JWT_SECRET")), map[string]interface{}{}))
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator with correct token secret but wrong alg passed want %v got %v",401, status)
	}

	//wrong claim
	userClaims := jwt.MapClaims{}
	userClaims["user"] = true
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(userClaims), nil); status != 401 || resp != "Unauthorized\n" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator with wrong claim passed want %v got %v",401, status)
	}

	// sending authorized requests
	adminClaims := jwt.MapClaims{}
	adminClaims["admin"] = true
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(adminClaims), nil); status != 200 || resp != "welcome" {
		t.Errorf(resp)
	} else {
		t.Logf("AdminAuthenticator with correct claim passed want %v got %v", 200, status)
	}
}
