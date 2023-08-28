package utils

import (
	"fmt"
	"time"
)

type Data struct {
	DateTime time.Time
	Number   int
}

func (d Data) String() string {
	return "number: " + fmt.Sprint(d.Number) + ", datetime: " + d.DateTime.String()
}
