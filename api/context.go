package api
import (
	"context"
	"net/http"
	"Golang_Project/pkg/model"


)
type contextKey string

// userContextKey is used as a key for getting and setting user information in the request
// context.
const userContextKey = contextKey("user")

// contextSetUser returns a new copy of the request with the provided User struct added to the
// context.
func (api *API) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context. The only time that
// this helper should be used is when we logically expect there to be a User struct value
// in the context, and if it doesn't exist it will firmly be an 'unexpected' error, upon we panic.
func (api *API) contextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
