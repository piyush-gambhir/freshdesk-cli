package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	os.Exit(cmd.Run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
