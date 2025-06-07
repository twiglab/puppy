package gbot

import "context"

type Commander interface {
	Execute(context.Context, any) error
}
