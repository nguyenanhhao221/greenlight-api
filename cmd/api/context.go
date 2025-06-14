package main

import (
	"context"
	"net/http"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
)

type contextKey string

// use this key as constant to get the user key from request context later
const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		// In practice, this function should always be called after we already set the user in the request context with the `authenticate` middleware
		// So this case will never happens, therefore it's better to panic here
		panic("panic: missing user value in request context")
	}
	return user
}
