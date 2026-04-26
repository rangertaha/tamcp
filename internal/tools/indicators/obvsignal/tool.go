// Package obvsignal registers OBV with EMA signal (utility).
package obvsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "obv_signal",
		Description: "On-Balance Volume with an EMA-smoothed signal line.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			c, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			o, sig := talib.OBVSIGNALFn(c, v, talib.ArgInt(args, "signal_period", 9))
			return talib.Two(o, sig, [2]string{"obv", "signal"}), talib.Tersum("obv_signal", o), nil
		},
	})
}
