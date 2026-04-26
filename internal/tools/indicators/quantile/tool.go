// Package quantile registers the rolling quantile indicator (Pandas TA).
package quantile

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "quantile",
		Description: "Rolling quantile q ∈ [0,1] with linear interpolation.",
		Group:       "statistic",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Input series"},
			{Name: "period", Type: "int", Default: 30, Desc: "rolling window"},
			{Name: "q", Type: "float", Default: 0.5, Desc: "quantile in [0,1]"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 30)
			if p <= 0 {
				p = 30
			}
			q := talib.ArgFloat(args, "q", 0.5)
			out := talib.QUANTILEFn(v, p, q)
			return talib.One(out), talib.Tersum("quantile", out), nil
		},
	})
}
