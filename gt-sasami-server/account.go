package gtsasamiserver

import (
	"time"
)

type Account struct {
	// ID (Auto generated)
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
}
