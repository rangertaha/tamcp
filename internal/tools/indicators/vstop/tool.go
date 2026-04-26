// Package vstop registers the Volatility Stop / Trailing ATR Stop (Pandas TA).
package vstop

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vstop",
		Description: "Volatility Stop: ATR-based trailing stop. Returns stop price and direction (+1/-1).",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 20, Desc: "ATR look-back"},
			{Name: "multiplier", Type: "float", Default: 2.0, Desc: "ATR offset multiplier"},
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
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			m := talib.ArgFloat(args, "multiplier", 2.0)
			if m <= 0 {
				m = 2.0
			}
			s, dir := talib.VSTOPFn(h, l, c, p, m)
			return talib.Two(s, dir, [2]string{"stop", "direction"}), talib.Tersum("vstop", s), nil
		},
	})
}
