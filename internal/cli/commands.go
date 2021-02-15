package cli

import (
    "errors"
    "fmt"
    "os"
    "strings"

    "github.com/leaanthony/clir"
    "github.com/olekukonko/tablewriter"

    "github.com/benoitgelineau/go-fetch-quotes/internal/utils"
)

func RegisterSearchAction(cli *clir.Cli) {
    search := cli.NewSubCommand("search", "Search a financial asset")
    search.LongDescription("Search a financial asset by name or ISIN and return the following information:\n\nSymbol | Name | Asset type | Last Price")

    search.Action(func() error {
        if len(os.Args) < 3 {
            return errors.New("Missing a value, please refer to the documentation by using `quotes search -help`")
        }

        // TODO handle string with space inside, encode URL
        query := search.OtherArgs()[0]

        fmt.Printf("Searching for %s...\n", query)
        assets := utils.ScrapeSearchResult(query)

        if len(assets) == 0 {
            fmt.Println("No result found.")
        } else {
            // Pretty print result in a table
            table := tablewriter.NewWriter(os.Stdout)
            table.SetHeader([]string{"Symbol", "Name", "Category", "Last price"})
            table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
            table.SetCenterSeparator("|")

            lines := [][]string{}
            for _, asset := range assets {
                line := []string{asset.Symbol, asset.Name, asset.Category, asset.LastPrice}
                lines = append(lines, line)
            }

            table.AppendBulk(lines)
            fmt.Println("Results found:\n")
            table.Render()
        }

        return nil
    })
}

func RegisterGetAction(cli *clir.Cli) {
    get := cli.NewSubCommand("get", "Return quotes")
}

func Get() error {
    return nil
}
