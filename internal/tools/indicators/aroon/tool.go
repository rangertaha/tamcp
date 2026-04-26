// Package aroon registers the aroon indicator with the talib dispatcher.
package aroon

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "aroon",
		Description: "Aroon up/down. Returns down, up.",
		Group:       "momentum",
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			d, u := talib.AROONFn(h, l, p)
			return talib.Two(d, u, [2]string{"down", "up"}), talib.Tersum("aroon", u), nil
		},
	})
}
