// Package cpr registers the Central Pivot Range (Pandas TA-style per-bar).
package cpr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pivot_cpr",
		Description: "Central Pivot Range using prior bar's HLC. Returns pivot, top central (tc), and bottom central (bc).",
		Group:       "overlap",
		Params:      talib.ParamsHLC(),
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
			p, tc, bc := talib.CPRFn(h, l, c)
			return talib.Three(p, tc, bc, [3]string{"pivot", "tc", "bc"}), talib.Tersum("pivot_cpr", p), nil
		},
	})
}
