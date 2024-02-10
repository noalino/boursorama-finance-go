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

type GetCommand struct{}

func Get() *cli.Command {
	var command Command[lib.GetResults] = GetCommand{}
	return &cli.Command{
		Name:      "get",
		Usage:     "Return historical data",
		UsageText: `bfinance get [options] SYMBOL`,
		Flags:     command.flags(),
		Action:    command.action,
	}
}

func (GetCommand) flags() []cli.Flag {
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

func (cmd GetCommand) action(cCtx *cli.Context) error {
	data, err := cmd.extract(cCtx)
	if err != nil {
		return err
	}

	cmd.load(cCtx, data)

	return nil
}

func (GetCommand) extract(cCtx *cli.Context) (lib.GetResults, error) {
	var symbol string

	if utils.IsDataFromPipe() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			symbol = strings.TrimSpace(s.Text())
			break
		}
	} else {
		if cCtx.NArg() == 0 {
			return lib.GetResults{}, errors.New("too few arguments, please refer to the documentation by using `bfinance get --help`")
		}
		symbol = cCtx.Args().First()
	}

	query := lib.GetQuery{
		Symbol:   symbol,
		From:     cCtx.String("from"),
		Duration: cCtx.String("duration"),
		Period:   cCtx.String("period"),
	}
	return lib.Get(query)
}

func (GetCommand) load(_ *cli.Context, data lib.GetResults) {
	if len(data) == 0 {
		fmt.Println("No data found.")
		return
	}

	fmt.Println("date,close,performance,high,low,open")
	for _, item := range data {
		fmt.Printf(
			"%s,%.2f,%s,%.2f,%.2f,%.2f\n",
			item.Date.Format(lib.GetResultDateFormat),
			item.Close,
			item.Perf,
			item.High,
			item.Low,
			item.Open,
		)
	}
}
