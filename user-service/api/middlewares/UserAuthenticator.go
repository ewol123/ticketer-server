package middlewares

import (
	"net/http"
	"github.com/go-chi/jwtauth"
)

func UserAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if claims["user"] != true {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if userIdStr, ok := claims["userId"].(string); ok {
			w.Header().Add("Requester-Id", userIdStr)
		}

		next.ServeHTTP(w, r)
	})
}
