// Package fcb registers Fractal Chaos Bands (Pandas TA legacy).
package fcb

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "fcb",
		Description: "Fractal Chaos Bands: most recent up-fractal high (upper) and down-fractal low (lower), carried forward.",
		Group:       "overlap",
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
			u, lo := talib.FCBFn(h, l)
			return talib.Two(u, lo, [2]string{"upper", "lower"}), talib.Tersum("fcb", u), nil
		},
	})
}
