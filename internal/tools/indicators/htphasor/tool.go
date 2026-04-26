// Package htphasor registers the ht_phasor indicator with the talib dispatcher.
package htphasor

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ht_phasor",
		Description: "Hilbert Transform Phasor (in_phase, quadrature).",
		Group:       "cycle",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			i, q := talib.HTPHASORFn(v)
			return talib.Two(i, q, [2]string{"in_phase", "quadrature"}), talib.Tersum("ht_phasor", i), nil
		},
	})
}
