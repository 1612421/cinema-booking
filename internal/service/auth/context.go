package authservice

import (
	"context"

	"go.uber.org/zap"
)

type ContextLogKey struct{}

type contextUserSessionKey struct{}

type ContextLogValue struct {
	APIName string
}

var (
	keyUserSession = contextUserSessionKey{}
)

func (c *ContextLogValue) ToLoggerFields() []zap.Field {
	if c == nil {
		return nil
	}
	return []zap.Field{
		zap.String("api_name", c.APIName),
	}
}

func GetUserSessionValue(ctx context.Context) interface{} {
	return ctx.Value(keyUserSession)
}

func NewContextWithUserSession(ctx context.Context, session interface{}) context.Context {
	return context.WithValue(ctx, keyUserSession, session)
}
