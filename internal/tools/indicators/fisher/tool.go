// Package fisher registers Ehlers' Fisher Transform (Pandas TA).
package fisher

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "fisher",
		Description: "Ehlers Fisher Transform of (H+L)/2. Returns fisher and a one-bar lagged signal series.",
		Group:       "momentum",
		Params:      talib.ParamsHLPeriod(9),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 9)
			if p <= 0 {
				p = 9
			}
			f, sig := talib.FISHERFn(h, l, p)
			return talib.Two(f, sig, [2]string{"fisher", "signal"}), talib.Tersum("fisher", f), nil
		},
	})
}
