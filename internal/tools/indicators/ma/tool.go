// Package ma registers the ma indicator with the talib dispatcher.
package ma

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ma",
		Description: "Moving Average; ma_type selects kernel.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 30)
			if p <= 0 {
				p = 30
			}
			mt, err := talib.MaTypeFromString(talib.ArgString(args, "matype", ""))
			if err != nil {
				return nil, "", err
			}
			out := talib.MA(v, p, mt)
			return talib.One(out), talib.Tersum("ma", out), nil
		},
	})
}
