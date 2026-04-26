// Package minmax registers the minmax indicator with the talib dispatcher.
package minmax

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "minmax",
		Description: "Rolling min and max. Returns min, max.",
		Group:       "operator",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 30)
			if p <= 0 {
				p = 30
			}
			mn, mx := talib.MINMAXFn(v, p)
			return talib.Two(mn, mx, [2]string{"min", "max"}), talib.Tersum("minmax", mn), nil
		},
	})
}
