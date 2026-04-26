// Package rocsignal registers ROC with an EMA signal line (utility).
package rocsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "roc_signal",
		Description: "Rate of Change with an EMA-smoothed signal line. Returns roc and signal.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 10, Desc: "ROC look-back"},
			{Name: "signal_period", Type: "int", Default: 9, Desc: "EMA period for signal"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			r, sig := talib.ROCSIGNALFn(v,
				talib.ArgInt(args, "period", 10),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(r, sig, [2]string{"roc", "signal"}), talib.Tersum("roc_signal", r), nil
		},
	})
}
