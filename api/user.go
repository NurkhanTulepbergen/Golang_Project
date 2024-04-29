package api

import (
	"Golang_Project/pkg/model"
	"Golang_Project/pkg/validator"
	"context"
	"errors"
	//"github.com/gorilla/mux"
	"net/http"
	"time"
)


func (api *API) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the request body into the anonymous struct.
	err := api.readJSON(w, r, &input)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if model.ValidateUser(v, user); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = api.UserModel.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			api.failedValidationResponse(w, r, v.Errors)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}
	//err = api.PermissionModel.AddForUser(user.ID, "shop:read")
	//if err != nil {
	//	api.serverErrorResponse(w, r, err)
	//	return
	//}

	err = api.PermissionModel.AddForUser(user.ID, "shop:read")
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := api.TokenModel.New(user.ID, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	api.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
}

func (api *API) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := api.readJSON(w, r, &input)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the details of the user associated with the token using the GetForToken() method.
	// If no matching record is found, then we let the client know that the token they provided
	// is not valid.
	user, err := api.UserModel.GetForToken(model.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			api.failedValidationResponse(w, r, v.Errors)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in the same
	// way that we did for our move records.
	err = api.UserModel.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			api.editConflictResponse(w, r)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully above, then delete all activation tokens for the user.
	err = api.TokenModel.DeleteAllForUser(model.ScopeActivation, user.ID)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	api.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
}

func (api *API) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the email and password from the request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := api.readJSON(w, r, &input)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	// Validate the email and password provided by the client.
	v := validator.New()
	model.ValidateEmail(v, input.Email)
	model.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Lookup the user record based on the email address. If no matching user was
	// found, then we call the app.invalidCredentialsResponse() helper to send a 401
	// Unauthorized response to the client (we will create this helper in a moment).
	user, err := api.UserModel.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			api.invalidCredentialsResponse(w, r)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}
	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		api.invalidCredentialsResponse(w, r)
		return
	}
	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := api.TokenModel.New(user.ID, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	token, err = api.TokenModel.New(user.ID, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	// Encode the token to JSON and send it in the response along with a 201 Created status code.
	err = api.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}
}
