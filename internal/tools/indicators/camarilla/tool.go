// Package camarilla registers Camarilla Pivot Points (per-bar from prior HLC).
package camarilla

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "camarilla",
		Description: "Camarilla Pivot Points using prior bar's HLC. Returns pp, r1..r4, s1..s4.",
		Group:       "overlap",
		Params:      talib.ParamsHLC(),
		Run: func(args map[string]any) (any, string, error) {
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
			pp, r1, r2, r3, r4, s1, s2, s3, s4 := talib.CAMARILLAFn(h, l, c)
			out := map[string]any{
				"pp": pp, "r1": r1, "r2": r2, "r3": r3, "r4": r4,
				"s1": s1, "s2": s2, "s3": s3, "s4": s4,
			}
			return out, talib.Tersum("camarilla", pp), nil
		},
	})
}
