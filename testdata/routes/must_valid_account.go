package routes

import (
	"context"

	"github.com/go-courier/courier"
)

type MustValidAccount struct {
	// Bearer access_token
	Authorization string `name:"Authorization,omitempty" in:"header"`
	// Bearer access_token in query
	AuthorizationInQuery string `name:"authorization,omitempty" in:"query"`
}

func (MustValidAccount) ContextKey() interface{} {
	return contextKeyAccount{}
}

func (c *MustValidAccount) Output(ctx context.Context) (interface{}, error) {
	return WithAccount(&Account{UserID: "admin"})(ctx), nil
}

type contextKeyAccount struct{}

func AccountFromContext(ctx context.Context) *Account {
	if a, ok := ctx.Value(contextKeyAccount{}).(*Account); ok {
		return a
	}
	return nil
}

func WithAccount(a *Account) courier.ContextWith {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextKeyAccount{}, a)
	}
}

type Account struct {
	UserID string
}
