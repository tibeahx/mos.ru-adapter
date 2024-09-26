package types

import "github.com/tibeahx/mos.ru-adapter/pkg/session"

type User struct {
	ID       string
	Email    string
	Password string
	Session  session.Session
}
