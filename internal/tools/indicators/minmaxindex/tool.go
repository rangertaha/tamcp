// Package minmaxindex registers the minmaxindex indicator with the talib dispatcher.
package minmaxindex

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "minmaxindex",
		Description: "Indices of rolling min and max.",
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
			mni, mxi := talib.MINMAXINDEXFn(v, p)
			return talib.Two(mni, mxi, [2]string{"min_idx", "max_idx"}), talib.Tersum("minmaxindex", mni), nil
		},
	})
}
