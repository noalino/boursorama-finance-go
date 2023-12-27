package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/noalino/boursorama-finance-go/internal/lib"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

type SearchCommand struct{}

func Search() *cli.Command {
	command := SearchCommand{}
	return &cli.Command{
		Name:      "search",
		Usage:     "Search for a financial asset",
		UsageText: "bfinance search [options] ASSET",
		Flags:     command.flags(),
		Action:    command.action,
	}
}

func (SearchCommand) flags() []cli.Flag {
	return []cli.Flag{
		&cli.UintFlag{
			Name:    "page",
			Aliases: []string{"P"},
			Value:   1,
			Usage:   "load specific page",
		},
		&cli.BoolFlag{
			Name:    "pretty",
			Aliases: []string{"p"},
			Value:   false,
			Usage:   "prettify the output",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   false,
			Usage:   "show more info",
		},
	}
}

func (SearchCommand) action(cCtx *cli.Context) error {
	if cCtx.NArg() == 0 {
		return errors.New("too few arguments, please refer to the documentation by running `bfinance search --help`")
	}

	query := lib.SearchQuery{Value: cCtx.Args().First(), Page: uint16(cCtx.Uint("page"))}

	utils.PrintfOrVoid(cCtx.Bool("verbose"), "Searching for '%s'...\n", query.Value)
	result, err := lib.Search(query)
	if err != nil {
		return err
	}

	if len(result.Assets) == 0 {
		fmt.Println("No result found.")
		return nil
	}

	fmt.Printf("Results found (page %d/%d):\n", result.Page, result.TotalPages)

	if cCtx.Bool("pretty") {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Symbol", "Name", "Market", "Close Price"})
		table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
		table.SetCenterSeparator("|")
		table.SetRowLine(true)

		lines := [][]string{}
		for _, asset := range result.Assets {
			line := []string{asset.Symbol, asset.Name, asset.Market, asset.LastPrice}
			lines = append(lines, line)
		}

		table.AppendBulk(lines)
		table.Render()
	} else {
		fmt.Println("symbol,name,market,last price")
		for _, asset := range result.Assets {
			fmt.Printf("%s,%s,%s,%s\n", asset.Symbol, asset.Name, asset.Market, asset.LastPrice)
		}
	}

	return nil
}
