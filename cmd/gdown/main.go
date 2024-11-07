package main

import (
	"github.com/haevg-rz/git-file-downloader/pkg/cli"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	"github.com/haevg-rz/git-file-downloader/pkg/log"
	"os"
)

func main() {
	if err := cli.Command().Execute(); err != nil {
		log.V(1).Printf("exit code: %d, error: %v\n", exit.Code, err)
		os.Exit(exit.Code)
	}
	os.Exit(exit.Success)
}
