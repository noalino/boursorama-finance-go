package options

import (
	"fmt"
	"time"
)

type From string

const (
	dateFormat = "02/01/2006"
)

func DefaultFrom() From {
	lastMonth := time.Now().AddDate(0, -1, 0)
	return From(lastMonth.Format(dateFormat))
}

func (f From) String() string {
	return string(f)
}

func (f From) ConvertToInternal() (string, error) {
	_, err := time.Parse(dateFormat, f.String())
	if err != nil {
		return "", fmt.Errorf("wrong date format: %v", err)
	}
	return f.String(), nil
}
