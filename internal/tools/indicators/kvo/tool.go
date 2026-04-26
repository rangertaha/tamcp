// Package kvo registers the Klinger Volume Oscillator (Pandas TA).
package kvo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kvo",
		Description: "Klinger Volume Oscillator: EMA(volume_force, fast) - EMA(volume_force, slow). Returns kvo and an EMA-smoothed signal.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "fast_period", Type: "int", Default: 34},
			{Name: "slow_period", Type: "int", Default: 55},
			{Name: "signal_period", Type: "int", Default: 13},
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
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			k, sig := talib.KVOFn(h, l, c, v,
				talib.ArgInt(args, "fast_period", 34),
				talib.ArgInt(args, "slow_period", 55),
				talib.ArgInt(args, "signal_period", 13),
			)
			return talib.Two(k, sig, [2]string{"kvo", "signal"}), talib.Tersum("kvo", k), nil
		},
	})
}
