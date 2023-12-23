package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/noalino/boursorama-finance-go/internal/lib"
	options "github.com/noalino/boursorama-finance-go/internal/lib/options/get"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

func Get() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Return quotes",
		UsageText: `quotes get [options] SYMBOL`,
		Flags:     initGetFlags(),
		Action:    getAction,
	}
}

func initGetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "duration",
			Aliases: []string{"d"},
			Value:   options.DefaultDuration.String(),
			Usage: fmt.Sprintf(`Specify the duration, it should be one of the following values:
	[%s]`, options.DurationsList),
		},
		&cli.StringFlag{
			Name:    "from",
			Aliases: []string{"f"},
			Value:   options.DefaultFrom().String(),
			Usage: `Specify the start date, it must be in the following format:
	DD/MM/YYYY`,
		},
		&cli.StringFlag{
			Name:    "period",
			Aliases: []string{"p"},
			Value:   options.DefaultPeriod.String(),
			Usage: fmt.Sprintf(`Specify the period, it should be one of the following values:
	[%s]`, options.PeriodsList),
		},
	}
}

func getAction(cCtx *cli.Context) error {
	var symbol string

	if utils.IsDataFromPipe() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			symbol = strings.TrimSpace(s.Text())
			break
		}
	} else {
		if cCtx.NArg() == 0 {
			return errors.New("too few arguments, please refer to the documentation by using `quotes get --help`")
		}
		symbol = cCtx.Args().First()
	}

	query := lib.GetQuery{
		Symbol:   symbol,
		From:     cCtx.String("from"),
		Duration: cCtx.String("duration"),
		Period:   cCtx.String("period"),
	}
	quotes, err := lib.Get(query)
	if err != nil {
		return err
	}

	if len(quotes) == 0 {
		fmt.Println("No quotes found.")
		return nil
	}

	fmt.Println("date,close,performance,high,low,open")
	for _, quote := range quotes {
		fmt.Printf(
			"%s,%.2f,%s,%.2f,%.2f,%.2f\n",
			quote.Date.Format(lib.GetResultDateFormat),
			quote.Close,
			quote.Perf,
			quote.High,
			quote.Low,
			quote.Open,
		)
	}

	return nil
}
