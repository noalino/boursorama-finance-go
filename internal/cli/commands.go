package cli

import (
    "errors"
    "fmt"
    "os"
    "strings"
    "time"

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
    // TODO
    get.LongDescription("Usage: quotes get [OPTIONS] [SYMBOL]")

    now := time.Now()
    lastMonth := now.AddDate(0,-1,0)
    // Default start date = a month from now
    startDate := lastMonth.Format(utils.LayoutISO)
    //var startDate string
    get.StringFlag("from", "Start date", &startDate)

    duration := "3M"
    get.StringFlag("duration", "Duration", &duration)

    period := "1"
    get.StringFlag("period", "Period", &period)

    get.Action(func() error {
        if len(os.Args) < 3 {
            return errors.New("Missing a value, please refer to the documentation by using `quotes get -help`")
        }
        // TODO Check flags
        symbol := strings.ToUpper(strings.TrimSpace(get.OtherArgs()[0]))

        startDateAsTime, err := time.Parse(utils.LayoutISO, startDate)
        if err != nil {
            return fmt.Errorf("Could not parse date: %v\n", err)
        }

        quotes := utils.GetQuotes(symbol, startDateAsTime, duration, period)
        if len(quotes) == 0 {
            fmt.Println("No quotes found.")
        } else {
            fmt.Printf("date,%s\n", symbol)
            for _, quote := range(quotes) {
                fmt.Printf("%s,%s\n", strings.TrimSpace(quote.Date), strings.TrimSpace(quote.Price))
            }
        }
        return nil
    })
}
