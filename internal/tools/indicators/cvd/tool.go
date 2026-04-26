// Package cvd registers a price-tick approximation of Cumulative Volume Delta.
package cvd

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cvd",
		Description: "Cumulative Volume Delta (close-tick approximation): cumulative sign(close - prev_close) * volume.",
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
			out := talib.CVDFn(c, v)
			return talib.One(out), talib.Tersum("cvd", out), nil
		},
	})
}
