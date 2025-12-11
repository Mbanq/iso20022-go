package common

import (
	"fmt"
	"time"
)

var EstLocation *time.Location

func init() {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone location: %v", err))
	}
	EstLocation = loc
}
