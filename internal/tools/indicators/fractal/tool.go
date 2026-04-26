// Package fractal registers Bill Williams Fractals (Pandas TA-style).
package fractal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "fractal",
		Description: "Bill Williams Fractals. Returns up and down (1.0 marks the centre bar i-2 of a 5-bar fractal pattern; 0 elsewhere).",
		Group:       "trend",
		Params:      talib.ParamsHL(),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			up, dn := talib.FRACTALFn(h, l)
			return talib.Two(up, dn, [2]string{"up", "down"}), talib.Tersum("fractal", up), nil
		},
	})
}
