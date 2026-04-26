// Package ssf registers Ehlers' 2-pole Super Smoother Filter (Pandas TA).
package ssf

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ssf",
		Description: "Ehlers Super Smoother Filter (2-pole Butterworth). Sharper roll-off than EMA at the same period.",
		Group:       "overlap",
		Params:      talib.ParamsRealPeriod(10),
		Run:         talib.RunRealPeriod("ssf", 10, talib.SSFFn),
	})
}
