package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	Password     string    `db:"password" json:"password,omitempty"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token,omitempty"`
	Role         Role      `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) AddRefreshToken(token string) {
	u.RefreshToken = token
	u.UpdatedAt = time.Now()
}

func (u *User) Validate() error {
	errorMessages := []string{}

	if len(u.Name) < 3 {
		errorMessages = append(errorMessages, "Name must be at least 3 characters long")
	}

	if len(u.Username) < 6 {
		errorMessages = append(errorMessages, "Username must be at least 6 characters long")
	}

	if len(u.Password) < 8 {
		errorMessages = append(errorMessages, "Password must be at least 8 characters long")
	}

	if len(u.Email) < 6 {
		errorMessages = append(errorMessages, "Email must be at least 6 characters long")
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func NewUser(name, username, email, password string) (*User, error) {
	user := &User{
		ID:        uuid.New(),
		Name:      name,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

type Role int

const (
	RoleUser Role = iota
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return "user"
	}
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
