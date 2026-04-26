// Package macdext registers the macdext indicator with the talib dispatcher.
package macdext

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "macdext",
		Description: "MACDEXT with controllable MA kinds.",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			fp := talib.ArgInt(args, "fast_period", 12)
			fma, err := talib.MaTypeFromString(talib.ArgString(args, "fast_matype", ""))
			if err != nil {
				return nil, "", err
			}
			sp := talib.ArgInt(args, "slow_period", 26)
			sma, err := talib.MaTypeFromString(talib.ArgString(args, "slow_matype", ""))
			if err != nil {
				return nil, "", err
			}
			gp := talib.ArgInt(args, "signal_period", 9)
			gma, err := talib.MaTypeFromString(talib.ArgString(args, "signal_matype", ""))
			if err != nil {
				return nil, "", err
			}
			m, sg, h := talib.MACDEXTFn(v, fp, fma, sp, sma, gp, gma)
			return talib.Three(m, sg, h, [3]string{"macd", "signal", "histogram"}), talib.Tersum("macdext", m), nil
		},
	})
}
