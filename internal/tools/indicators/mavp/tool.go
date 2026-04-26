// Package mavp registers the mavp indicator with the talib dispatcher.
package mavp

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mavp",
		Description: "Moving Average with Variable Period.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			periods, err := talib.ArgFloats(args, "periods")
			if err != nil {
				return nil, "", err
			}
			mn := talib.ArgInt(args, "min_period", 2)
			mx := talib.ArgInt(args, "max_period", 30)
			mt, err := talib.MaTypeFromString(talib.ArgString(args, "matype", ""))
			if err != nil {
				return nil, "", err
			}
			out := talib.MAVPFn(v, periods, mn, mx, mt)
			return talib.One(out), talib.Tersum("mavp", out), nil
		},
	})
}
