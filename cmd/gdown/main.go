package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/haevg-rz/git-file-downloader/pkg/cli"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
)

func gracefulExit() {
	cli.Done <- true
	log.FileLogWg.Wait()
	os.Exit(exit.Code.Int())
}

func main() {
	sigCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func(ctx context.Context) {
		<-ctx.Done()
		log.V(1).Printf("exit code: %d, error: received signal interrupt", exit.ReceivedSignal)
		exit.Code = exit.ReceivedSignal
		gracefulExit()
	}(sigCtx)

	if err := cli.Command().Execute(); err != nil {
		log.V(1).Printf("exit code: %d, error: %v\n", exit.Code, err)
	} else {
		log.V(1).Printf("exit code: %d\n", exit.Code)
	}

	gracefulExit()
}
