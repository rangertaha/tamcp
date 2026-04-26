// Package kvopct registers KVO normalised by its rolling absolute mean (utility).
package kvopct

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kvo_pct",
		Description: "KVO normalised: 100 * KVO / SMA(|KVO|, slow_period). Comparable across instruments.",
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
			out := talib.KVOPCTFn(h, l, c, v,
				talib.ArgInt(args, "fast_period", 34),
				talib.ArgInt(args, "slow_period", 55),
				talib.ArgInt(args, "signal_period", 13),
			)
			return talib.One(out), talib.Tersum("kvo_pct", out), nil
		},
	})
}
