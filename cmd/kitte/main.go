package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/yuhaiin/kitte/internal/update"
)

func main() {
	log.SetFlags(0)

	root := flag.String("root", ".", "repository root")
	flag.Parse()

	if flag.NArg() != 1 || flag.Arg(0) != "update" {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] update\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := update.Run(ctx, *root); err != nil {
		log.Fatal(err)
	}
}
