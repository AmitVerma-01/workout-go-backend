package middleware

import (
	"context"
	"fmt"
	"go_beginner/internals/store"
	"go_beginner/internals/tokens"
	"go_beginner/utils"
	"net/http"
	"strings"
)


type UserMiddleware struct {
	UserStore store.UserStore 
}

type ContextKey string

const UserContextKey ContextKey = "user"

func SetUser(r *http.Request, user *store.User) *http.Request{
	ctx := context.WithValue(r.Context(), UserContextKey , user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User{
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		// panic("could not get user from context")
		return store.AnonymousUser
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"Error": "Invalid or expired token"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"Error": "Invalid or expired token"})
			return
		}

		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"Error": "Invalid or expired token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"Error": "You must be logged in to access this resource"})
			return
		}
		fmt.Printf("User ID in RequireUser middleware: %+v\n", user)
		next.ServeHTTP(w, r)
	})
}
