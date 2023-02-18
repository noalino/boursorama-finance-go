package main

import (
    "fmt"

    "github.com/leaanthony/clir"

    commands "internal/cli"
)

func main() {
    cli := clir.NewCli("Quotes", "A basic scraper tool to get financial assets quotes", "v0.1.0")

    commands.RegisterSearchAction(cli)
    commands.RegisterGetAction(cli)

    if err := cli.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}

