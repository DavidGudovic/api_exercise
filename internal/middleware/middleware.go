package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/DavidGudovic/api_exercise/internal/store"
	"github.com/DavidGudovic/api_exercise/internal/tokens"
	"github.com/DavidGudovic/api_exercise/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)

	if !ok {
		panic("Missing user in request context")
	}

	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid authorization header"})
			return
		}

		tokenString := headerParts[1]

		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, tokenString)

		if err != nil {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid or expired token"})
			return
		}

		if user == nil {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid or expired token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return
	})
}

func (um *UserMiddleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "You must be authenticated to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
