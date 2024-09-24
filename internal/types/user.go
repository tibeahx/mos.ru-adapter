package types

import "test-task/pkg/session"

type User struct {
	ID       string
	Email    string
	Password string
	Session  session.Session
}
