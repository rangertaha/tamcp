// Package bbpsignal registers BBP with an EMA signal line (utility).
package bbpsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bbp_signal",
		Description: "Bollinger %B with an EMA-smoothed signal line.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 20},
			{Name: "nbdevup", Type: "float", Default: 2.0},
			{Name: "nbdevdn", Type: "float", Default: 2.0},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			b, sig := talib.BBPSIGNALFn(v,
				talib.ArgInt(args, "period", 20),
				talib.ArgFloat(args, "nbdevup", 2.0),
				talib.ArgFloat(args, "nbdevdn", 2.0),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(b, sig, [2]string{"bbp", "signal"}), talib.Tersum("bbp_signal", b), nil
		},
	})
}
