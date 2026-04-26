// Package woodie registers Woodie Pivot Points (per-bar from prior HLC).
package woodie

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "woodie",
		Description: "Woodie Pivot Points using prior bar's HLC. Returns pp, r1, r2, s1, s2.",
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
			pp, r1, r2, s1, s2 := talib.WOODIEFn(h, l, c)
			out := map[string]any{"pp": pp, "r1": r1, "r2": r2, "s1": s1, "s2": s2}
			return out, talib.Tersum("woodie", pp), nil
		},
	})
}
