// Package htsine registers the ht_sine indicator with the talib dispatcher.
package htsine

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ht_sine",
		Description: "Hilbert Transform SineWave (sine, lead_sine).",
		Group:       "cycle",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			s, ls := talib.HTSINEFn(v)
			return talib.Two(s, ls, [2]string{"sine", "lead_sine"}), talib.Tersum("ht_sine", s), nil
		},
	})
}
