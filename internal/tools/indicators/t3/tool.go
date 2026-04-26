// Package t3 registers the t3 indicator with the talib dispatcher.
package t3

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "t3",
		Description: "Tillson T3.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 5)
			if p <= 0 {
				p = 5
			}
			vf := talib.ArgFloat(args, "v_factor", 0.7)
			out := talib.T3Fn(v, p, vf)
			return talib.One(out), talib.Tersum("t3", out), nil
		},
	})
}
