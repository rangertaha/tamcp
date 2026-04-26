// Package donchian registers Donchian Channels (Pandas TA, cinar).
package donchian

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "donchian",
		Description: "Donchian Channels. Returns upper, middle, lower over the look-back period.",
		Group:       "overlap",
		Params:      talib.ParamsHLPeriod(20),
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			u, m, lo := talib.DONCHIANFn(h, l, p)
			return talib.Three(u, m, lo, [3]string{"upper", "middle", "lower"}), talib.Tersum("donchian", m), nil
		},
	})
}
