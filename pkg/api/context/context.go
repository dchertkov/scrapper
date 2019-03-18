package context

import (
	"context"

	"github.com/dchertkov/scrapper/pkg/types"
	"github.com/sirupsen/logrus"
)

type key int

const (
	logKey key = iota
	serviceKey
)

func FromLog(ctx context.Context) *logrus.Entry {
	return ctx.Value(logKey).(*logrus.Entry)
}

func ToLog(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, logKey, log)
}

func FromService(ctx context.Context) (*types.Service, bool) {
	s, ok := ctx.Value(serviceKey).(*types.Service)
	return s, ok
}

func ToService(ctx context.Context, s *types.Service) context.Context {
	return context.WithValue(ctx, serviceKey, s)
}
