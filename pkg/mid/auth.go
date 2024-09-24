package mid

import "net/http"

type AuthMiddleware struct{}

func (m *AuthMiddleware) Authenthicate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
