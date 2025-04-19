package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/karan/watchlist/services"
)
type contextKey string

const UserIDKey contextKey = "userId"
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc( func (w http.ResponseWriter , r *http.Request){
		authHeader := r.Header.Get("Authorization")
		if authHeader =="" || !strings.HasPrefix(authHeader , "Bearer"){
			http.Error(w,"Missing or invalid token",http.StatusUnauthorized)
			return
		}
		tokenString:= strings.TrimPrefix(authHeader , "Bearer ");
		userId , err := services.VerifyToken(tokenString)
		if err!=nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx:= context.WithValue(r.Context(),UserIDKey,userId)
		next.ServeHTTP(w , r.WithContext(ctx))
		/*
		1. ctx := context.WithValue(r.Context(), userIDKey, userId)
This line creates a new context (ctx) that contains the userId value, and it associates this value with a key (userIDKey) within the context.

r.Context() gets the current context of the incoming HTTP request r. Every HTTP request in Go automatically has a context associated with it. The context carries values, deadlines, cancellation signals, and other request-specific data.

context.WithValue(parentContext, key, value) creates a new context that is derived from the parentContext (in this case, r.Context()). It associates the value (userId) with a specific key (userIDKey) in the new context. This does not modify the original context (r.Context()), but rather creates a new one that includes the additional data (userId).

userIDKey is the key that identifies the value in the context (usually defined as a custom type to avoid conflicts with other keys, as shown in earlier examples).

userId is the value that you want to associate with the key in the context (this could be the ID of a user, for example).

The result is that ctx now contains all the data from r.Context() plus the added userId.

2. next.ServeHTTP(w, r.WithContext(ctx))
This line calls the next handler in the HTTP middleware chain and passes the new context (ctx) along with the request.

next represents the next handler in the HTTP request-response lifecycle. In Go's HTTP server, this is typically the next middleware or the final handler that will process the request.

r.WithContext(ctx) creates a new request object (r), but it associates it with the new context (ctx). The original request (r) is not modified; a new instance of the request is created with the updated context. This new request is then passed along to the next handler.

r.WithContext(ctx) is used to make sure the request that reaches the next handler has the updated context (which now contains the user ID).

next.ServeHTTP(w, r.WithContext(ctx)) forwards the new request (with the updated context) to the next handler, allowing that handler to access the user ID (userId) through the context.
*/ 
	})
}