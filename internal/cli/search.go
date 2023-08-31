package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/leaanthony/clir"
	"github.com/olekukonko/tablewriter"

	"github.com/noalino/boursorama-finance-go/internal/utils"
)

func RegisterSearchAction(cli *clir.Cli) {
	search := cli.NewSubCommand("search", "Search a financial asset")
	search.LongDescription(
		`
Search a financial asset by name or ISIN and return the following information:
Symbol, Name, Category, Last price

Usage: quotes search [name | ISIN]`)

	// Flags
	var pretty, verbose bool
	search.BoolFlag("pretty", "Display output in a table.", &pretty)
	search.BoolFlag("verbose", "Log more info.", &verbose)

	// Actions
	search.Action(func() error {
		otherArgs := search.OtherArgs()
		if len(otherArgs) == 0 {
			return errors.New("too few arguments, please refer to the documentation by using `quotes search -help`")
		}

		query := utils.SearchQuery{Value: otherArgs[0]}

		utils.PrintfOrVoid(verbose, "Searching for '%s'...\n", query)
		assets, err := utils.ScrapeSearchResult(query)
		if err != nil {
			return err
		}

		if len(assets) == 0 {
			fmt.Println("No result found.")
			return nil
		}

		utils.PrintlnOrVoid(verbose, "Results found:")

		if pretty {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Symbol", "Name", "Market", "Last price"})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetRowLine(true)

			lines := [][]string{}
			for _, asset := range assets {
				line := []string{asset.Symbol, asset.Name, asset.Market, asset.LastPrice}
				lines = append(lines, line)
			}

			table.AppendBulk(lines)
			table.Render()
		} else {
			fmt.Println("symbol,name,market,last price")
			for _, asset := range assets {
				fmt.Printf("%s,%s,%s,%s\n", asset.Symbol, asset.Name, asset.Market, asset.LastPrice)
			}
		}

		return nil
	})
}
