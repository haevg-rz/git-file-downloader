package main

import (
	"github.com/haevg-rz/git-file-downloader/pkg/cli"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func gracefulExit() {
	cli.Done <- true
	os.Exit(exit.Code)
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.V(1).Printf("exit code: %d, error: received signal '%v'", exit.ReceivedSignal, <-sig)
		exit.Code = exit.ReceivedSignal
		gracefulExit()
	}()

	if err := cli.Command().Execute(); err != nil {
		log.V(1).Printf("exit code: %d, error: %v\n", exit.Code, err)
	} else {
		log.V(1).Printf("exit code: %d\n", exit.Code)
	}

	gracefulExit()
}
