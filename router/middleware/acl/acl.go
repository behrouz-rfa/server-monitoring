package acl

import (

	"net/http"
	"server-monitoring/shared/session"
)

// DisallowAuth does not allow authenticated users to access the page
func DisallowAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is authenticated, don't allow them to access the page
		if sess.Values["id"] != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	}
}

// DisallowAnon does not allow anonymous users to access the page
func DisallowAnon(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	}
}
