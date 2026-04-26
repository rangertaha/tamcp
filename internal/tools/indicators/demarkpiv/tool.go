// Package demarkpiv registers DeMark Pivot Points (per-bar from prior OHLC).
package demarkpiv

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "demark_pivots",
		Description: "DeMark Pivot Points using prior bar's OHLC. Returns pp, r1, s1.",
		Group:       "overlap",
		Params:      talib.ParamsOHLC(),
		Run: func(args map[string]any) (any, string, error) {
			o, err := talib.ArgFloats(args, "open")
			if err != nil {
				return nil, "", err
			}
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
			pp, r1, s1 := talib.DEMARKFn(o, h, l, c)
			return talib.Three(pp, r1, s1, [3]string{"pp", "r1", "s1"}), talib.Tersum("demark_pivots", pp), nil
		},
	})
}
