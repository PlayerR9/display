package screen

import "context"

type contextKey struct{}

type Context struct {
}

func NewContext() context.Context {
	c := &Context{}

	ctx := context.WithValue(context.Background(), contextKey{}, c)

	return ctx
}
