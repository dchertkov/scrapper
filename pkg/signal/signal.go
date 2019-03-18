package signal

import (
	"context"
	"os"
	"os/signal"
)

func WithCancel(ctx context.Context, fn func(sig os.Signal)) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)

		select {
		case sig := <-ch:
			fn(sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
