// Package momsignal registers Momentum with EMA signal (utility).
package momsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mom_signal",
		Description: "Momentum with an EMA-smoothed signal line.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 10},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			m, sig := talib.MOMSIGNALFn(v,
				talib.ArgInt(args, "period", 10),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(m, sig, [2]string{"mom", "signal"}), talib.Tersum("mom_signal", m), nil
		},
	})
}
