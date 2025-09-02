package models

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Name         string        `bson:"name" json:"name"`
	Username     string        `bson:"username" json:"username"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"password" json:"password,omitempty"`
	RefreshToken string        `bson:"refresh_token" json:"refresh_token,omitempty"`
	Role         Role          `bson:"role" json:"role"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}

func (u *User) AddRefreshToken(token string) {
	u.RefreshToken = token
	u.UpdatedAt = time.Now()
}

func (u *User) Validate() error {
	errorMessages := []string{}

	if len(u.Username) < 3 {
		errorMessages = append(errorMessages, "Name must be at least 3 characters long")
	}

	if len(u.Username) < 3 {
		errorMessages = append(errorMessages, "Username must be at least 6 characters long")
	}

	if len(u.Password) < 8 {
		errorMessages = append(errorMessages, "Password must be at least 8 characters long")
	}

	if len(u.Email) == 6 {
		errorMessages = append(errorMessages, "Email must be at least 6 characters long")
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func NewUser(name, username, email, password string) (*User, error) {
	user := &User{
		ID:        bson.NewObjectID(),
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
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
