// Package pvt registers the Price Volume Trend (Pandas TA, cinar).
package pvt

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pvt",
		Description: "Price Volume Trend (cumulative): pvt[i] = pvt[i-1] + volume[i] * (close[i] - close[i-1]) / close[i-1].",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
		},
		Run: func(args map[string]any) (any, string, error) {
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			out := talib.PVTFn(c, v)
			return talib.One(out), talib.Tersum("pvt", out), nil
		},
	})
}
