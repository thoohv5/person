package context

import (
	"context"
	"errors"
)

type Context struct {
	ILogicalMessage
	context.Context
}

func WithLM(ctx context.Context) *Context {
	return &Context{
		ILogicalMessage: NewComptroller(),
		Context:         ctx,
	}
}

const (
	ErrNotFoundKey = "not found key"
)

type ILogicalMessage interface {
	Key() string
	Marshal(ctx context.Context) (string, error)
	Unmarshal(data string) error
	Merge(message ILogicalMessage)
}

var (
	_logicalMessage = make(map[string]ILogicalMessage)
)

func registerLogicalMessage(msg ILogicalMessage) {
	if _, ok := _logicalMessage[msg.Key()]; !ok {
		_logicalMessage[msg.Key()] = msg
	}
}

func GetLogicalMessage() map[string]ILogicalMessage {
	return _logicalMessage
}

func WithMessage(ctx context.Context, msg ILogicalMessage) (context.Context, error) {
	for key, message := range GetLogicalMessage() {
		if key == msg.Key() {
			if err := GetMessage(ctx, message); err != nil {
				if ErrNotFoundKey == err.Error() {
					err = nil
					continue
				}
				return nil, err
			}
			msg.Merge(message)
		}
	}
	s, err := msg.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, msg.Key(), s), nil
}

func GetMessage(ctx context.Context, msg ILogicalMessage) error {
	s, ok := ctx.Value(msg.Key()).(string)
	if !ok {
		return errors.New(ErrNotFoundKey)
	}
	return msg.Unmarshal(s)
}
