// Package vo registers the Volume Oscillator (Pandas TA, cinar).
package vo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vo",
		Description: "Volume Oscillator: 100 * (EMA(vol, fast) - EMA(vol, slow)) / EMA(vol, slow).",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume series"},
			{Name: "fast_period", Type: "int", Default: 14},
			{Name: "slow_period", Type: "int", Default: 28},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			out := talib.VOFn(v,
				talib.ArgInt(args, "fast_period", 14),
				talib.ArgInt(args, "slow_period", 28),
			)
			return talib.One(out), talib.Tersum("vo", out), nil
		},
	})
}
