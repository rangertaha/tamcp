// Package vortex registers the Vortex Indicator (Pandas TA).
package vortex

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vortex",
		Description: "Vortex Indicator: VI+ and VI− over the look-back period.",
		Group:       "momentum",
		Params:      talib.ParamsHLCPeriod(14),
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
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			vip, vim := talib.VORTEXFn(h, l, c, p)
			return talib.Two(vip, vim, [2]string{"vi_plus", "vi_minus"}), talib.Tersum("vortex", vip), nil
		},
	})
}
