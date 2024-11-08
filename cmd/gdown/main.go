package main

import (
	"github.com/haevg-rz/git-file-downloader/pkg/cli"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.V(1).Printf("exit code: %d, error: received signal '%v'", exit.ReceivedSignal, <-sig)
		os.Exit(exit.ReceivedSignal)
	}()

	if err := cli.Command().Execute(); err != nil {
		log.V(1).Printf("exit code: %d, error: %v\n", exit.Code, err)
		os.Exit(exit.Code)
	}
	log.V(1).Println("done")
	os.Exit(exit.Success)
}
