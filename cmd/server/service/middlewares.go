package service

import (
	"context"

	"github.com/lichuan0620/secret-keeper-backend/internal/queueclient"
	servicemodel "github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
)

var contextKeyQueueClient interface{} = new(byte)

func WithQueueClient(qc queueclient.Interface) servicemodel.Middleware {
	return func(ctx context.Context, f func(context.Context)) {
		f(context.WithValue(ctx, contextKeyQueueClient, qc))
	}
}

func GetQueueClient(ctx context.Context) queueclient.Interface {
	return ctx.Value(contextKeyQueueClient).(queueclient.Interface)
}
