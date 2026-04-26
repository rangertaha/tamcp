// Package thermo registers Elder's Thermometer (Pandas TA).
package thermo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "thermo",
		Description: "Elder Thermometer: max(|Δhigh|, |Δlow|) and its EMA-smoothed value.",
		Group:       "volatility",
		Params:      talib.ParamsHLPeriod(20),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			t, sm := talib.THERMOFn(h, l, p)
			return talib.Two(t, sm, [2]string{"thermo", "smoothed"}), talib.Tersum("thermo", sm), nil
		},
	})
}
