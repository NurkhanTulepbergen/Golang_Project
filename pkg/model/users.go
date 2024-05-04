package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"Golang_Project/pkg/validator"
)

var AnonymousUser = &User{}

// User represents an individual user.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// IsAnonymous checks if the user is anonymous.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// password represents a user's password.
type password struct {
	plaintext *string
	hash      []byte
}

// Set calculates the bcrypt hash of a plaintext password.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// Matches checks whether the provided plaintext password matches the hashed password.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// UserModel wraps the connection pool.
type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// ErrDuplicateEmail represents a duplicate email error.
var ErrDuplicateEmail = errors.New("duplicate email")

// Insert a new record for the user.
func (m *UserModel) Insert(user *User) error {
	query := `
INSERT INTO users (name, email, password, activated, version)
VALUES ($1, $2, $3, $4, 1)
RETURNING id, created_at, version + 1`
	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// GetByEmail retrieves the user details from the database based on the email address.
func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
SELECT id, created_at, name, email, password, activated, version
FROM users
WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update the details for a specific user.
func (m *UserModel) Update(user *User) error {
	query := `
UPDATE users
SET name = $1, email = $2, password = $3, activated = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		return err
	}
	return nil
}

// ValidateEmail validates the email format.
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext validates the plaintext password.
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// ValidateUser validates the user data.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// GetForToken retrieves the user based on the token.
func (m *UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
SELECT users.id, users.created_at, users.name, users.email, users.password, users.activated, users.version
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Token represents a user token.
