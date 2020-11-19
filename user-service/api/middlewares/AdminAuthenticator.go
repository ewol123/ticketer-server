package middlewares

import (
	"github.com/go-chi/jwtauth"
	"net/http"
)

func AdminAuthenticator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if claims["admin"] != true {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if userIdStr, ok := claims["userId"].(string); ok {
			w.Header().Add("Requester-Id", userIdStr)
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
