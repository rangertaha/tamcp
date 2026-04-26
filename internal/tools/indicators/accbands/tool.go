// Package accbands registers the Acceleration Bands indicator (Pandas TA, cinar).
package accbands

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "accbands",
		Description: "Acceleration Bands. Returns upper, middle, lower SMA over `period` bars.",
		Group:       "volatility",
		Params:      talib.ParamsHLCPeriod(20),
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
			p := talib.ArgInt(args, "period", 20)
			if p <= 0 {
				p = 20
			}
			u, m, lo := talib.ACCBANDSFn(h, l, c, p)
			return talib.Three(u, m, lo, [3]string{"upper", "middle", "lower"}), talib.Tersum("accbands", m), nil
		},
	})
}
