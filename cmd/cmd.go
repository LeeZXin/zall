package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func initWaitContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(
			signalChannel,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		select {
		case <-signalChannel:
		case <-ctx.Done():
		}
		cancel()
		signal.Reset()
	}()
	return ctx, cancel
}
