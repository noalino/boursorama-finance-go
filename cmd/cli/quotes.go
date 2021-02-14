package main

import (
    "errors"
    "log"
    "os"

    "github.com/leaanthony/clir"

    commands "github.com/benoitgelineau/go-fetch-quotes/internal/cli"
)

func main() {
    cli := clir.NewCli("Quotes", "A basic wrapper to get financial asset quotes", "v0.0.1")

    search := cli.NewSubCommand("search", "Search a financial asset")
    search.LongDescription("Search a financial asset by name or ISIN and return the following information:\n\nSymbol | Name | Asset type | Last Price")
    search.Action(func() error {
        return handleAction(commands.Search)
    })

    get := cli.NewSubCommand("get", "Return quotes")
    get.Action(func() error {
        return handleAction(commands.Get)
    })

    err := cli.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func handleAction(executeAction commands.Command) error {
    if len(os.Args) < 3 {
        return errors.New("Missing a value, please refer to the documentation by using `quotes -help`")
    }
    return executeAction()
}
