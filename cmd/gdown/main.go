package main

import (
	"github.com/haevg-rz/git-file-downloader/pkg/cli"
	"github.com/pkg/errors"
	"log"
)

func main() {
	if err := cli.Command().Execute(); err != nil {
		log.Fatal(errors.Wrap(err, "error while executing command"))
	}
}
