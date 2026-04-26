// Package bbands registers the bbands indicator with the talib dispatcher.
package bbands

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bbands",
		Description: "Bollinger Bands. Returns upper, middle, lower.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 5)
			if p <= 0 {
				p = 5
			}
			du := talib.ArgFloat(args, "nbdevup", 2)
			dd := talib.ArgFloat(args, "nbdevdn", 2)
			mt, err := talib.MaTypeFromString(talib.ArgString(args, "matype", ""))
			if err != nil {
				return nil, "", err
			}
			u, m, l := talib.BBANDSFn(v, p, du, dd, mt)
			return talib.Three(u, m, l, [3]string{"upper", "middle", "lower"}), talib.Tersum("bbands", m), nil
		},
	})
}
