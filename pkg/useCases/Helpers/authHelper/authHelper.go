package authHelper

import (
	"fmt"
	"net/http"
	"os"
)

// BasicAuth implements a simple middleware handler for adding basic http auth to a route.
func BasicAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w)
				return
			}
			if user != "lila" || pass != os.Getenv("BACK_PASS") {
				w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, "Lila"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, "Lila"))
	w.WriteHeader(http.StatusUnauthorized)
}
