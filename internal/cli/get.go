package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/noalino/boursorama-finance-go/internal/options"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

func (cli *Cli) RegisterGetAction() {
	get := cli.NewSubCommand("get", "Return quotes\n")
	get.LongDescription("Usage: quotes get [OPTIONS] SYMBOL")

	// Flags
	from := options.DefaultFrom().String()
	duration := options.DefaultDuration.String()
	period := options.DefaultPeriod.String()

	get.StringFlag("from",
		`Specify the start date, it must be in the following format:
DD/MM/YYYY`,
		&from)

	get.StringFlag("duration",
		`Specify the duration, it should be one of the following values:
[`+options.DurationsList.String()+`]`, &duration)

	get.StringFlag("period",
		`Specify the period, it should be one the following values:
[`+options.PeriodsList.String()+`]`, &period)

	// Actions
	get.Action(func() error {

		var symbol string

		if utils.IsDataFromPipe() {
			s := bufio.NewScanner(os.Stdin)
			for s.Scan() {
				symbol = strings.TrimSpace(s.Text())
				break
			}
		} else {
			otherArgs := get.OtherArgs()
			if len(otherArgs) == 0 {
				return errors.New("too few arguments, please refer to the documentation by using `quotes get -help`")
			}
			symbol = otherArgs[0]
		}

		query := utils.QuotesQuery{
			Symbol:   symbol,
			From:     from,
			Duration: duration,
			Period:   period,
		}
		quotes, err := utils.GetQuotes(query)
		if err != nil {
			return err
		}

		if len(quotes) == 0 {
			fmt.Println("No quotes found.")
			return nil
		}

		fmt.Printf("date,%s\n", symbol)
		for _, quote := range quotes {
			fmt.Printf("%s,%.2f\n", quote.Date, quote.Price)
		}

		return nil
	})
}
