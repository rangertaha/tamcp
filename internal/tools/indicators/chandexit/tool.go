// Package chandexit registers the Chandelier Exit indicator (Pandas TA).
package chandexit

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "chandelier_exit",
		Description: "Chandelier Exit: ATR-offset trailing stops from rolling high (long) and rolling low (short).",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 22, Desc: "look-back for HHV/LLV and ATR"},
			{Name: "multiplier", Type: "float", Default: 3.0, Desc: "ATR offset multiplier"},
		},
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
			p := talib.ArgInt(args, "period", 22)
			if p <= 0 {
				p = 22
			}
			m := talib.ArgFloat(args, "multiplier", 3.0)
			if m <= 0 {
				m = 3.0
			}
			lo, sh := talib.CHANDELIEREXITFn(h, l, c, p, m)
			return talib.Two(lo, sh, [2]string{"long_exit", "short_exit"}), talib.Tersum("chandelier_exit", lo), nil
		},
	})
}
