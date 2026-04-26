// Package volratio registers the Volume Ratio indicator (cinar).
package volratio

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vol_ratio",
		Description: "Volume Ratio: volume[i] / SMA(volume, period). >1 means above-average volume.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Volume series"},
			{Name: "period", Type: "int", Default: 20, Desc: "SMA window"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			out := talib.VOLRATIOFn(v, p)
			return talib.One(out), talib.Tersum("vol_ratio", out), nil
		},
	})
}
