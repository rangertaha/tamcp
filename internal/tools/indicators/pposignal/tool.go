// Package pposignal registers PPO with an EMA signal line (utility).
package pposignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ppo_signal",
		Description: "Percentage Price Oscillator with an EMA-smoothed signal line.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "fast_period", Type: "int", Default: 12},
			{Name: "slow_period", Type: "int", Default: 26},
			{Name: "signal_period", Type: "int", Default: 9},
			{Name: "matype", Type: "string", Default: "EMA", Desc: "moving-average kernel for PPO"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			mt, err := talib.MaTypeFromString(talib.ArgString(args, "matype", "EMA"))
			if err != nil {
				return nil, "", err
			}
			p, sig := talib.PPOSIGNALFn(v,
				talib.ArgInt(args, "fast_period", 12),
				talib.ArgInt(args, "slow_period", 26),
				talib.ArgInt(args, "signal_period", 9),
				mt,
			)
			return talib.Two(p, sig, [2]string{"ppo", "signal"}), talib.Tersum("ppo_signal", p), nil
		},
	})
}
