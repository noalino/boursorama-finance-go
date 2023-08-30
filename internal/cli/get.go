package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/leaanthony/clir"

	"github.com/noalino/boursorama-finance-go/internal/utils"
)

func RegisterGetAction(cli *clir.Cli) {
	get := cli.NewSubCommand("get", "Return quotes")
	get.LongDescription(
		`
Usage: quotes get [OPTIONS] SYMBOL`)

	lastMonth := time.Now().AddDate(0, -1, 0)
	startDate := lastMonth.Format(utils.LayoutISO)
	defaultStartDate := "a month from now"

	// Flags
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

		validSymbol := utils.ValidateInput(symbol)
		if validSymbol == "" {
			return errors.New("symbol value must be valid and not empty")
		}

		startDateAsTime, err := time.Parse(utils.LayoutISO, startDate)
		if err != nil {
			return fmt.Errorf("wrong date format: %v", err)
		}

		quotes, err := utils.GetQuotes(validSymbol, startDateAsTime, duration, period)
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
