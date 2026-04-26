// Package supertrend registers the SuperTrend indicator (Pandas TA, cinar).
package supertrend

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "supertrend",
		Description: "SuperTrend: ATR-based trend follower. Returns the active band as `trend` and `direction` ∈ {+1,-1}.",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 10, Desc: "ATR look-back"},
			{Name: "multiplier", Type: "float", Default: 3.0, Desc: "ATR band multiplier"},
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
			p := talib.ArgInt(args, "period", 10)
			if p <= 0 {
				p = 10
			}
			m := talib.ArgFloat(args, "multiplier", 3.0)
			if m <= 0 {
				m = 3.0
			}
			tr, dir := talib.SUPERTRENDFn(h, l, c, p, m)
			return talib.Two(tr, dir, [2]string{"trend", "direction"}), talib.Tersum("supertrend", tr), nil
		},
	})
}
