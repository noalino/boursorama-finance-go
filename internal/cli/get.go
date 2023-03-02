package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/leaanthony/clir"

	"internal/utils"
)

func RegisterGetAction(cli *clir.Cli) {
	get := cli.NewSubCommand("get", "Return quotes")
	get.LongDescription(
		`
Usage: quotes get [OPTIONS] SYMBOL`)

	now := time.Now()
	// Default start date = a month from now
	lastMonth := now.AddDate(0, -1, 0)
	startDate := lastMonth.Format(utils.LayoutISO)
	defaultStartDate := "a month from now"
	get.StringFlag("from",
		`Specify the start date, it must be in the following format:
DD/MM/YYYY`,
		&defaultStartDate)

	duration := utils.DefaultDurations[2]
	get.StringFlag("duration",
		`Specify the duration, it should be one of the following values:
[`+strings.Join(utils.DefaultDurations, ", ")+`]`, &duration)

	period := utils.DefaultPeriods[0]
	get.StringFlag("period",
		`Specify the period, it should be one the following values:
[`+strings.Join(utils.DefaultPeriods, ", ")+`]`, &period)

	get.Action(func() error {
		otherArgs := get.OtherArgs()
		if len(otherArgs) == 0 {
			return errors.New("Too few arguments, please refer to the documentation by using `quotes get -help`")
		}

		symbol := otherArgs[0]
		validSymbol := utils.ValidateInput(symbol)
		if validSymbol == "" {
			return errors.New("Symbol value must be valid and not empty.")
		}

		startDateAsTime, err := time.Parse(utils.LayoutISO, startDate)
		if err != nil {
			return fmt.Errorf("Wrong date format: %v\n", err)
		}

		quotes, err := utils.GetQuotes(validSymbol, startDateAsTime, duration, period)
		if err != nil {
			return err
		}

		if len(quotes) == 0 {
			fmt.Println("No quotes found.")
		} else {
			fmt.Printf("date,%s\n", symbol)
			for _, quote := range quotes {
				fmt.Printf("%s,%.2f\n", quote.Date, quote.Price)
			}
		}
		return nil
	})
}
