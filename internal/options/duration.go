package options

import (
	"fmt"
	"strings"
)

type Duration string
type Durations []Duration

const (
	oneMonth     Duration = "1M"
	twoMonths    Duration = "2M"
	threeMonths  Duration = "3M"
	fourMonths   Duration = "4M"
	fiveMonths   Duration = "5M"
	sixMonths    Duration = "6M"
	sevenMonths  Duration = "7M"
	eightMonths  Duration = "8M"
	nineMonths   Duration = "9M"
	tenMonths    Duration = "10M"
	elevenMonths Duration = "11M"
	twelveMonths Duration = "12M"
	oneYear      Duration = "1Y"
	twoYears     Duration = "2Y"
	threeYears   Duration = "3Y"
)

var DefaultDuration = threeMonths
var DurationsList = Durations{oneMonth, twoMonths, threeMonths, fourMonths, fiveMonths, sixMonths, sevenMonths, eightMonths, nineMonths, tenMonths, elevenMonths, twelveMonths, oneYear, twoYears, threeYears}

func (d Duration) String() string {
	return string(d)
}

func (d Duration) ConvertToInternal() (string, error) {
	for _, duration := range DurationsList {
		if duration == d {
			return string(d), nil
		}
	}
	return "", fmt.Errorf("duration must be one of %v", DurationsList)
}

func (d Durations) String() string {
	durationsStr := make([]string, len(d))
	for i := range d {
		durationsStr[i] = string(d[i])
	}
	return strings.Join(durationsStr, ", ")
}
