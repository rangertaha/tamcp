// Package apo registers the apo indicator with the talib dispatcher.
package apo

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "apo",
		Description: "Absolute Price Oscillator",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			fp := talib.ArgInt(args, "fast_period", 12)
			sp := talib.ArgInt(args, "slow_period", 26)
			mt, err := talib.MaTypeFromString(talib.ArgString(args, "matype", ""))
			if err != nil {
				return nil, "", err
			}
			out := talib.APOFn(v, fp, sp, mt)
			return talib.One(out), talib.Tersum("apo", out), nil
		},
	})
}
