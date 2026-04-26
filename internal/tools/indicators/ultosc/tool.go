// Package ultosc registers the ultosc indicator with the talib dispatcher.
package ultosc

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ultosc",
		Description: "Ultimate Oscillator",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			p1 := talib.ArgInt(args, "period1", 7)
			p2 := talib.ArgInt(args, "period2", 14)
			p3 := talib.ArgInt(args, "period3", 28)
			out := talib.ULTOSCFn(h, l, c, p1, p2, p3)
			return talib.One(out), talib.Tersum("ultosc", out), nil
		},
	})
}
