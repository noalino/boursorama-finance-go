package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/noalino/boursorama-finance-go/internal/lib"
	options "github.com/noalino/boursorama-finance-go/internal/lib/options/get"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

type getFlags struct {
	duration string
	from     string
	period   string
}

func (cli *Cli) RegisterGetAction() {
	get := cli.NewSubCommand("get", "Return quotes\n")
	get.LongDescription("Usage: quotes get [OPTIONS] SYMBOL")

	// Flags
	flags := &getFlags{
		duration: options.DefaultDuration.String(),
		from:     options.DefaultFrom().String(),
		period:   options.DefaultPeriod.String(),
	}

	get.StringFlag(
		"from",
		"Specify the start date, it must be in the following format:\nDD/MM/YYYY",
		&flags.from)

	get.StringFlag(
		"duration",
		fmt.Sprintf("Specify the duration, it should be one of the following values:\n[%s]", options.DurationsList),
		&flags.duration)

	get.StringFlag(
		"period",
		fmt.Sprintf("Specify the period, it should be one of the following values:\n[%s]", options.PeriodsList),
		&flags.period)

	// Action
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

		query := lib.GetQuery{
			Symbol:   symbol,
			From:     flags.from,
			Duration: flags.duration,
			Period:   flags.period,
		}
		quotes, err := lib.Get(query)
		if err != nil {
			return err
		}

		if len(quotes) == 0 {
			fmt.Println("No quotes found.")
			return nil
		}

		fmt.Printf("date,%s\n", symbol)
		for _, quote := range quotes {
			fmt.Printf("%s,%.2f\n", quote.Date.Format(lib.GetResultDateFormat), quote.Price)
		}

		return nil
	})
}
