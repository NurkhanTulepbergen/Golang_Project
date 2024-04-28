// middleware.go

package api

import (
	"Golang_Project/pkg/model"
	"Golang_Project/pkg/validator"
	"errors"
	"net/http"
	"strings"
)

// Middleware functions need access to API methods, so embed the API struct within them.

// authenticate middleware checks the authentication token.
func (api *API) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any
		// caches that the response may vary based on the value of the Authorization
		// header in the request.
		w.Header().Add("Vary", "Authorization")
		// Retrieve the value of the Authorization header from the request. This will
		// return the empty string "" if there is no such header found.
		authorizationHeader := r.Header.Get("Authorization")
		// If there is no Authorization header found, use the contextSetUser() helper
		// that we just made to add the AnonymousUser to the request context. Then we
		// call the next handler in the chain and return without executing any of the
		// code below.
		if authorizationHeader == "" {
			r = api.contextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		// Otherwise, we expect the value of the Authorization header to be in the format
		// "Bearer <token>". We try to split this into its constituent parts, and if the
		// header isn't in the expected format we return a 401 Unauthorized response
		// using the invalidAuthenticationTokenResponse() helper (which we will create
		// in a moment).
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			api.invalidAuthenticationTokenResponse(w, r)
			return
		}
		// Extract the actual authentication token from the header parts.
		token := headerParts[1]
		// Validate the token to make sure it is in a sensible format.
		v := validator.New()
		// If the token isn't valid, use the invalidAuthenticationTokenResponse()
		// helper to send a response, rather than the failedValidationResponse() helper
		// that we'd normally use.
		if model.ValidateTokenPlaintext(v, token); !v.Valid() {
			api.invalidAuthenticationTokenResponse(w, r)
			return
		}
		// Retrieve the details of the user associated with the authentication token,
		// again calling the invalidAuthenticationTokenResponse() helper if no
		// matching record was found. IMPORTANT: Notice that we are using
		// ScopeAuthentication as the first parameter here.
		user, err := api.UserModel.GetForToken(model.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				api.invalidAuthenticationTokenResponse(w, r)
			default:
				api.serverErrorResponse(w, r, err)
			}
			return
		}
		// Call the contextSetUser() helper to add the user information to the request
		// context.
		r = api.contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})

}

func (api *API) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use the contextGetUser helper to retrieve the user information from the request context.
		user := api.contextGetUser(r)

		// If the user is anonymous, then call authenticationRequiredResponse to inform the client
		// that they should be authenticated before trying again.
		if user.IsAnonymous() {
			api.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}
func (api *API) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the contextGetUser() helper that we made earlier to retrieve the user
		// information from the request context.
		user := api.contextGetUser(r)
		// If the user is anonymous, then call the authenticationRequiredResponse() to
		// inform the client that they should authenticate before trying again.
		if user.IsAnonymous() {
			api.authenticationRequiredResponse(w, r)
			return
		}
		// If the user is not activated, use the inactiveAccountResponse() helper to
		// inform them that they need to activate their account.
		if !user.Activated {
			api.inactiveAccountResponse(w, r)
			return
		}
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (api *API) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the request context.
		user := api.contextGetUser(r)
		// Get the slice of permissions for the user.
		permissions, err := api.PermissionModel.GetAllForUser(user.ID)
		if err != nil {
			api.serverErrorResponse(w, r, err)
			return
		}
		// Check if the slice includes the required permission. If it doesn't, then
		// return a 403 Forbidden response.
		if !permissions.Include(code) {
			api.notPermittedResponse(w, r)
			return
		}
		// Otherwise they have the required permission so we call the next handler in
		// the chain.
		next.ServeHTTP(w, r)
	}
	// Wrap this with the requireActivatedUser() middleware before returning it.
	return api.requireActivatedUser(fn)
}
