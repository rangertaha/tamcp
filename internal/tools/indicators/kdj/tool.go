// Package kdj registers the KDJ stochastic indicator (cinar, Yatala).
package kdj

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "kdj",
		Description: "KDJ stochastic. Returns k, d, j series. K and D are 2/3-1/3 smoothings of the raw stochastic; J = 3K - 2D.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 9, Desc: "stochastic look-back"},
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
			p := talib.ArgInt(args, "period", 9)
			if p <= 0 {
				p = 9
			}
			k, d, j := talib.KDJFn(h, l, c, p)
			return talib.Three(k, d, j, [3]string{"k", "d", "j"}), talib.Tersum("kdj", k), nil
		},
	})
}
