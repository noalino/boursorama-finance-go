package cli

import (
	"errors"
	"fmt"
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
	lastMonth := now.AddDate(0, -1, 0)
	// Default start date = a month from now
	startDate := lastMonth.Format(utils.LayoutISO)
	defaultStartDate := "a month from now"
	get.StringFlag("from",
		`Specify the start date, it must be in the following format:
DD/MM/YYYY`,
		&defaultStartDate)

	duration := "3M"
	get.StringFlag("duration",
		`Specify the duration, it should be one of the following values:
["1M","2M","3M","4M","5M","6M","7M","8M","9M","10M","11M","1Y","2Y","3Y"]`,
		&duration)

	period := "1"
	get.StringFlag("period",
		`Specify the period, it should be one the following values:
["1","7","30","365"]`,
		&period)

	get.Action(func() error {
		otherArgs := get.OtherArgs()
		if len(otherArgs) == 0 {
			return errors.New("Too few arguments, please refer to the documentation by using `quotes get -help`")
		}

		symbol := otherArgs[0]

		startDateAsTime, err := time.Parse(utils.LayoutISO, startDate)
		if err != nil {
			return fmt.Errorf("Wrong date format: %v\n", err)
		}

		quotes, err := utils.GetQuotes(symbol, startDateAsTime, duration, period)
		if err != nil {
			return err
		}

		if len(quotes) == 0 {
			fmt.Println("No quotes found.")
		} else {
			fmt.Printf("date,%s\n", symbol)
			for _, quote := range quotes {
				fmt.Printf("%s,%s\n", quote.Date, quote.Price)
			}
		}
		return nil
	})
}
