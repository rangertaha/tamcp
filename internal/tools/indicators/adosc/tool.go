// Package adosc registers the adosc indicator with the talib dispatcher.
package adosc

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "adosc",
		Description: "Chaikin A/D Oscillator",
		Group:       "volume",
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
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			fp := talib.ArgInt(args, "fast_period", 3)
			sp := talib.ArgInt(args, "slow_period", 10)
			out := talib.ADOSCFn(h, l, c, v, fp, sp)
			return talib.One(out), talib.Tersum("adosc", out), nil
		},
	})
}
