package main

import (
	"fmt"
	"os"

	"github.com/leaanthony/clir"

	. "github.com/noalino/boursorama-finance-go/internal/cli"
)

func main() {
	cli := Cli{
		Cli: clir.NewCli(
			"quotes",
			"A basic scraper tool to get financial assets quotes",
			"v1.4.0",
		),
	}

	cli.Init()

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
