package cli

import (
    "fmt"
    "os"
    "strings"

    "github.com/olekukonko/tablewriter"

    "github.com/benoitgelineau/go-fetch-quotes/internal/utils"
)

type Command func() error

func Search() error {
    query := os.Args[2]

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
}

func Get() error {
    return nil
}
