// Package mama registers the mama indicator with the talib dispatcher.
package mama

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mama",
		Description: "MESA Adaptive Moving Average. Returns mama, fama.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			fl := talib.ArgFloat(args, "fast_limit", 0.5)
			sl := talib.ArgFloat(args, "slow_limit", 0.05)
			m, f := talib.MAMAFn(v, fl, sl)
			return talib.Two(m, f, [2]string{"mama", "fama"}), talib.Tersum("mama", m), nil
		},
	})
}
