package main

import (
	"fmt"
	"os"

	"github.com/leaanthony/clir"

	commands "github.com/noalino/boursorama-finance-go/internal/cli"
)

func main() {
	cli := clir.NewCli("Quotes", "A basic scraper tool to get financial assets quotes", "v1.2.0")

	commands.RegisterSearchAction(cli)
	commands.RegisterGetAction(cli)

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
