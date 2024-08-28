package http

import "context"

type ClientInterface interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
}
