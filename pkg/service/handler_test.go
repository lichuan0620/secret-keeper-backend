package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
)

func TestMiddlewareOrder(t *testing.T) {
	var buf []int
	middlewareLn := parseMiddlewares([]model.Middleware{
		func(ctx context.Context, f func(context.Context)) {
			buf = append(buf, 1)
			f(ctx)
		},
		func(ctx context.Context, f func(context.Context)) {
			buf = append(buf, 2)
			f(ctx)
		},
		func(ctx context.Context, f func(context.Context)) {
			f(ctx)
			buf = append(buf, 4)
		},
	})
	middlewareLn.execute(context.Background(), func(ctx context.Context) {
		buf = append(buf, 3)
	})
	if expected := []int{1, 2, 3, 4}; !reflect.DeepEqual(buf, expected) {
		t.Fatalf("expecting %v; got %v", expected, buf)
	}
}

func TestMiddlewareContext(t *testing.T) {
	type keyT string
	const key keyT = "k"

	parseMiddlewares([]model.Middleware{
		func(ctx context.Context, f func(context.Context)) {
			f(context.WithValue(ctx, key, "v1"))
		},
		func(ctx context.Context, f func(context.Context)) {
			f(context.WithValue(ctx, key, "v2"))
		},
	}).execute(context.Background(), func(ctx context.Context) {
		if got := ctx.Value(key).(string); got != "v2" {
			t.Errorf(`expecting %v; got %v`, "v2", got)
		}
	})

	parseMiddlewares([]model.Middleware{
		func(ctx context.Context, f func(context.Context)) {
			f(context.WithValue(ctx, key, "v"))
		},
		func(_ context.Context, f func(context.Context)) {
			f(context.Background())
		},
	}).execute(context.Background(), func(ctx context.Context) {
		if got := ctx.Value(key); got != nil {
			t.Errorf("expecting nil, got %v", got)
		}
	})
}
