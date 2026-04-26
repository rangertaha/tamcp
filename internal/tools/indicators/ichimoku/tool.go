// Package ichimoku registers the Ichimoku Cloud (Pandas TA, cinar).
package ichimoku

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ichimoku",
		Description: "Ichimoku Cloud. Returns tenkan, kijun, senkou_a, senkou_b, chikou. Per-bar element = what's plotted at that bar (Pandas TA convention).",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "tenkan", Type: "int", Default: 9},
			{Name: "kijun", Type: "int", Default: 26},
			{Name: "senkou", Type: "int", Default: 52},
			{Name: "displacement", Type: "int", Default: 26},
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
			tk, kj, sa, sb, ch := talib.ICHIMOKUFn(h, l, c,
				talib.ArgInt(args, "tenkan", 9),
				talib.ArgInt(args, "kijun", 26),
				talib.ArgInt(args, "senkou", 52),
				talib.ArgInt(args, "displacement", 26),
			)
			out := map[string]any{
				"tenkan":   tk,
				"kijun":    kj,
				"senkou_a": sa,
				"senkou_b": sb,
				"chikou":   ch,
			}
			return out, talib.Tersum("ichimoku", kj), nil
		},
	})
}
