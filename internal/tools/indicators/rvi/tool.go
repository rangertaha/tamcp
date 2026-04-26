// Package rvi registers the Relative Vigor Index (Pandas TA).
package rvi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "rvi",
		Description: "Relative Vigor Index. Returns rvi and a 4-bar SWMA signal.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "open", Type: "number[]", Required: true, Desc: "Open prices"},
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 10, Desc: "rolling sum window"},
		},
		Run: func(args map[string]any) (any, string, error) {
			o, err := talib.ArgFloats(args, "open")
			if err != nil {
				return nil, "", err
			}
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
			p := talib.ArgInt(args, "period", 10)
			if p <= 0 {
				p = 10
			}
			r, sig := talib.RVIFn(o, h, l, c, p)
			return talib.Two(r, sig, [2]string{"rvi", "signal"}), talib.Tersum("rvi", r), nil
		},
	})
}
