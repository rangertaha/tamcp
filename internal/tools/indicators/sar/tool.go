// Package sar registers the sar indicator with the talib dispatcher.
package sar

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sar",
		Description: "Parabolic SAR.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			acc := talib.ArgFloat(args, "acceleration", 0.02)
			mx := talib.ArgFloat(args, "maximum", 0.2)
			out := talib.SARFn(h, l, acc, mx)
			return talib.One(out), talib.Tersum("sar", out), nil
		},
	})
}
