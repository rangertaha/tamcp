// Package obvsmooth registers an EMA-smoothed On-Balance Volume (cinar).
package obvsmooth

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "obv_smoothed",
		Description: "EMA-smoothed On-Balance Volume.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 21, Desc: "EMA period"},
		},
		Run: func(args map[string]any) (any, string, error) {
			c, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 21)
			if p <= 0 {
				p = 21
			}
			out := talib.OBVSMOOTHFn(c, v, p)
			return talib.One(out), talib.Tersum("obv_smoothed", out), nil
		},
	})
}
