// Package stddev registers the stddev indicator with the talib dispatcher.
package stddev

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "stddev",
		Description: "Standard deviation",
		Group:       "statistic",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 5)
			if p <= 0 {
				p = 5
			}
			nd := talib.ArgFloat(args, "nbdev", 1)
			out := talib.STDDEVFn(v, p, nd)
			return talib.One(out), talib.Tersum("stddev", out), nil
		},
	})
}
