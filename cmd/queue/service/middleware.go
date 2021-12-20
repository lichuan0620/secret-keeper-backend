package service

import (
	"context"

	"github.com/lichuan0620/secret-keeper-backend/internal/queue"
	servicemodel "github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
)

var contextKeyQueue interface{} = new(byte)

func WithQueue(q queue.Interface) servicemodel.Middleware {
	return func(ctx context.Context, f func(context.Context)) {
		f(context.WithValue(ctx, contextKeyQueue, q))
	}
}

func GetQueue(ctx context.Context) queue.Interface {
	return ctx.Value(contextKeyQueue).(queue.Interface)
}
