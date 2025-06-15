package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
	"github.com/nguyenanhhao221/greenlight-api/internal/models"
	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a defer function which go will always run in the event of a panic as Go unwinds the stack
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimitMiddleware(next http.Handler) http.Handler {
	// Declare variable here will make the return function closure , it kept the reference to these variable
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// launch a background go routine that run every minute and check if last seen is more than 3 minutes and perform clean up
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			// Split to get the IP address for ip base rate limit
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Prevent concurrent access to the map
			mu.Lock()
			_, found := clients[ip]
			if !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any
		// caches that the response may vary based on the value of the Authorization
		// header in the request.
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.logger.Info("Authorization header not found, setting user as AnonymousUser")
			// If header is not set, treat user as anonymous user
			app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		// Validate if token header correctly set in form as "Bearer <token>"
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 {
			app.invalidTokenResponse(w, r)
			return
		}
		if headerParts[0] != "Bearer" {
			app.invalidTokenResponse(w, r)
			return
		}
		token := headerParts[1]
		v := validator.New()
		if data.ValidatePlaintextToken(v, token); !v.Valid() {
			app.failValidationResponse(w, r, v.Errors)
			return
		}

		// Validate that token actual relate to a user in database
		user, err := app.models.User.GetUserWithToken(token, data.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.invalidTokenResponse(w, r)
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		// Now that a user and token are valid, set the user in request context
		app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
