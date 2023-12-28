package main

import (
	"fmt"
	"os"

	commands "github.com/noalino/boursorama-finance-go/internal/cli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "bfinance",
		Version: "v2.0.0",
		Usage:   "A basic scraper tool to get financial assets historical data",
		Commands: []*cli.Command{
			commands.Get(),
			commands.Search(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
