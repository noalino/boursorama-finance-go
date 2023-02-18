package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/leaanthony/clir"
	"github.com/olekukonko/tablewriter"

	"internal/utils"
)

func RegisterSearchAction(cli *clir.Cli) {
	search := cli.NewSubCommand("search", "Search a financial asset")
	search.LongDescription(
		`
Search a financial asset by name or ISIN and return the following information:
-----------------------------------------
| Symbol | Name | Category | Last Price |
-----------------------------------------

Usage: quotes search NAME | ISIN`)

	search.Action(func() error {
		otherArgs := search.OtherArgs()
		if len(otherArgs) == 0 {
			return errors.New("Too few arguments, please refer to the documentation by using `quotes search -help`")
		}

		query := otherArgs[0]
		validQuery := utils.ValidateInput(query)
		if validQuery == "" {
			return errors.New("Search value must be valid and not empty.")
		}

		fmt.Printf("Searching for '%s'...\n", validQuery)
		assets, err := utils.ScrapeSearchResult(validQuery)
		if err != nil {
			return err
		}

		if len(assets) == 0 {
			fmt.Println("No result found.")
		} else {
			// Pretty print result in a table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Symbol", "Name", "Last price"})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			table.SetRowLine(true)

			lines := [][]string{}
			for _, asset := range assets {
				line := []string{asset.Symbol, asset.Name, asset.LastPrice}
				lines = append(lines, line)
			}

			table.AppendBulk(lines)
			fmt.Println("Results found:\n")
			table.Render()
		}

		return nil
	})
}
