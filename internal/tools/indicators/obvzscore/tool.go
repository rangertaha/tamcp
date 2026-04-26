// Package obvzscore registers Z-score of OBV (utility).
package obvzscore

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "obv_zscore",
		Description: "Rolling Z-score of On-Balance Volume.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 30},
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
			p := talib.ArgInt(args, "period", 30)
			if p <= 0 {
				p = 30
			}
			out := talib.OBVZSCOREFn(c, v, p)
			return talib.One(out), talib.Tersum("obv_zscore", out), nil
		},
	})
}
