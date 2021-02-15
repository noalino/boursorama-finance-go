package main

import (
    "fmt"

    "github.com/leaanthony/clir"

    commands "github.com/benoitgelineau/go-fetch-quotes/internal/cli"
)

func main() {
    cli := clir.NewCli("Quotes", "A basic wrapper to get financial asset quotes", "v0.0.1")

    commands.RegisterSearchAction(cli)
    commands.RegisterGetAction(cli)

    if err := cli.Run(); err != nil {
        fmt.Printf("Error while trying to run the CLI:\n%v\n", err)
    }
}

