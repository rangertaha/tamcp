// Package wt registers LazyBear's Wave Trend Oscillator (Pandas TA).
package wt

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "wt",
		Description: "LazyBear Wave Trend Oscillator. Returns wt1 and wt2 (SMA(wt1, 4)).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "n1", Type: "int", Default: 10, Desc: "EMA period for ESA and the deviation"},
			{Name: "n2", Type: "int", Default: 21, Desc: "EMA period for the channel index"},
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
			w1, w2 := talib.WTFn(h, l, c,
				talib.ArgInt(args, "n1", 10),
				talib.ArgInt(args, "n2", 21),
			)
			return talib.Two(w1, w2, [2]string{"wt1", "wt2"}), talib.Tersum("wt", w1), nil
		},
	})
}
