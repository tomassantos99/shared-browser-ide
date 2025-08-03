package middleware

import "net/http"

func SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// TODO: perform session checks

		next.ServeHTTP(w,r)
	})
}
