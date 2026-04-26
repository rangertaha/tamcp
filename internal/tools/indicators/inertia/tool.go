// Package inertia registers Pandas TA's Inertia indicator.
package inertia

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "inertia",
		Description: "Inertia: linear regression smoothing of RVI over `regression_period`.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "open", Type: "number[]", Required: true, Desc: "Open prices"},
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "rvi_period", Type: "int", Default: 14},
			{Name: "regression_period", Type: "int", Default: 20},
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
			rp := talib.ArgInt(args, "rvi_period", 14)
			gp := talib.ArgInt(args, "regression_period", 20)
			out := talib.INERTIAFn(o, h, l, c, rp, gp)
			return talib.One(out), talib.Tersum("inertia", out), nil
		},
	})
}
