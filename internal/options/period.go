package options

import (
	"fmt"
	"strings"
)

type Period string
type Periods []Period

const (
	daily   Period = "1"
	weekly  Period = "7"
	monthly Period = "30"
	yearly  Period = "365"
)

var DefaultPeriod = daily
var PeriodsList = Periods{daily, weekly, monthly, yearly}

var periodMap = map[Period]string{
	daily:   "1",
	weekly:  "7",
	monthly: "30",
	yearly:  "365",
}

func (p Period) String() string {
	return string(p)
}

func (p Period) ConvertToInternal() (string, error) {
	if value, found := periodMap[p]; !found {
		return "", fmt.Errorf("period must be one of %v", PeriodsList)
	} else {
		return value, nil
	}
}

func (p Periods) String() string {
	periodsStr := make([]string, len(p))
	for i := range p {
		periodsStr[i] = string(p[i])
	}
	return strings.Join(periodsStr, ", ")
}
