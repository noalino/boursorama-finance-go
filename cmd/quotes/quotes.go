package main

import (
	"fmt"

	"github.com/leaanthony/clir"

	commands "github.com/noalino/boursorama-finance-go/internal/cli"
)

func main() {
	cli := clir.NewCli("Quotes", "A basic scraper tool to get financial assets quotes", "v1.1.1")

	commands.RegisterSearchAction(cli)
	commands.RegisterGetAction(cli)

	if err := cli.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
